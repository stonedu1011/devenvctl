package plan

import (
	"fmt"
	dockerclient "github.com/docker/docker/client"
	"github.com/stonedu1011/devenvctl/pkg/devenv"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func NewScriptHookExecutables(hook devenv.Hook, wd string, vars []string, searchDirs ...string) ([]Executable, error) {
	// note: we assume the value of hook is a script filename in resource directory
	str, ok := hook.Value.(string)
	if !ok {
		return nil, fmt.Errorf(`expected script hook to have string value, but got %v`, hook.Value)
	}
	searchPaths := make([]string, 0, len(searchDirs) * 2)
	for _, dir := range searchDirs {
		searchPaths = append(searchPaths,
			filepath.Join(dir, string(hook.Phase), str),
			filepath.Join(dir, string(hook.Phase)+"-"+str),
		)
	}

	var cmd string
	for _, path := range searchPaths {
		if stat, e := os.Stat(path); e == nil && !stat.IsDir() {
			cmd = path
			break
		}
	}
	if len(cmd) == 0 {
		return nil, fmt.Errorf(`hook script [%s] not found in [%s]`, str, strings.Join(searchDirs, ", "))
	}

	exec := &ShellExecutable{
		Cmds: []string{cmd},
		WD:   wd,
		Env:  vars,
		Desc: fmt.Sprintf(`%v shell`, hook.Phase),
	}
	return []Executable{exec}, nil
}

// NewContainerHookExecutables create an executable that monitor existing containers and wait for it to stop
// Note: Currently, we don't support manually start containers. So the hook containers should be started.
//       Therefore, only post-start containers are possible.
func NewContainerHookExecutables(dockerClient *dockerclient.Client, phase devenv.HookPhase,
	_ DockerComposePlanMetadata, cResolver ContainerResolver, hooks ...devenv.Hook) ([]Executable, error) {

	if phase != devenv.PhasePostStart {
		return nil, fmt.Errorf(`container hooks are only supported in post-start phase`)
	}

	containers := make([]string, 0, len(hooks))
	execs := make([]Executable, 0, 2)
	for i := range hooks {
		if hooks[i].Type != devenv.TypeContainer {
			continue
		}
		name, ok := hooks[i].Value.(string)
		if !ok {
			return nil, fmt.Errorf(`expected hook value to be string, but got %v`, hooks[i].Value)
		}
		containers = append(containers, name)
	}
	// TODO support dynamic timeout
	exec := NewContainerMonitorExecutable(dockerClient, func(exec *ContainerMonitorExecutable) {
		exec.Names = containers
		exec.Desc = fmt.Sprintf(`%v containers`, phase)
		exec.Resolver = cResolver
	})
	execs = append(execs, exec.WithTimeout(30 * time.Second))
	return execs, nil
}
