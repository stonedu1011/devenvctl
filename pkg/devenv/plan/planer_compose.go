package plan

import (
	"context"
	"fmt"
	lanaiutils "github.com/cisco-open/go-lanai/pkg/utils"
	"github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	"github.com/stonedu1011/devenvctl/pkg/devenv"
	"github.com/stonedu1011/devenvctl/pkg/utils"
	"github.com/stonedu1011/devenvctl/pkg/utils/tmplutils"
	"os"
	"path/filepath"
	"time"
)

const defaultComposeFile = `docker-compose.yml`

type DockerComposePlanMetadata struct {
	Profile       *devenv.Profile
	Variables     []devenv.Variable
	WorkingDir    string
	DockerVersion types.Version
	ComposePath   string
	ResourceDir   string
}

func NewDockerComposePlanner(p *devenv.Profile, wd string) *DockerComposePlanner {
	plan := DockerComposePlanner{
		Profile:    p,
		WorkingDir: utils.AbsPath(wd, p.FS),
	}
	return &plan
}

type DockerComposePlanner struct {
	// WorkingDir the working directory. Usually is the temporary dir configured by rootcmd.GlobalArgs
	WorkingDir   string
	Profile      *devenv.Profile
	metadata     DockerComposePlanMetadata
	dockerClient *dockerclient.Client
}

func (pl *DockerComposePlanner) Prepare() (err error) {
	logger.Infof(`Using Docker Compose`)
	defer func() {
		if err == nil {
			logger.Infof(`Working directory is ready: %s`, pl.WorkingDir)
		}
	}()

	// generate metadata
	pl.metadata = DockerComposePlanMetadata{
		Profile:    pl.Profile,
		WorkingDir: pl.WorkingDir,
	}

	// load compose template
	tmplPath := pl.Profile.ComposePath()
	logger.Debugf(`Loading [%s]`, tmplPath)
	tmpl, e := tmplutils.NewTemplate().ParseFS(pl.Profile.FS, tmplPath)
	if e != nil {
		return fmt.Errorf("unable to process docker compose template [%s]: %v", utils.AbsPath(tmplPath, pl.Profile.FS), e)
	}

	// create a docker-compose.yml
	pl.metadata.ComposePath = filepath.Join(pl.WorkingDir, defaultComposeFile)
	logger.Debugf(`Finalizing docker compose file: %s`, pl.metadata.ComposePath)
	switch fi, e := os.Stat(pl.WorkingDir); {
	case e != nil:
		return fmt.Errorf(`unable to access directory [%s]: %v`, pl.WorkingDir, e)
	case !fi.IsDir():
		return fmt.Errorf(`unable to access directory [%s]: not a directory`, pl.WorkingDir)
	}
	composeF, e := os.OpenFile(pl.metadata.ComposePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if e != nil {
		return fmt.Errorf(`unable to generate docker compose [%s]: %v`, pl.metadata.ComposePath, e)
	}
	defer func() { _ = composeF.Close() }()

	// generate docker-compose.yml
	if e := tmpl.ExecuteTemplate(composeF, filepath.Base(tmplPath), pl); e != nil {
		return fmt.Errorf(`unable to generate docker compose [%s]: %v`, pl.metadata.ComposePath, e)
	}

	// copy resources
	srcResPath := pl.Profile.ResourceDir()
	pl.metadata.ResourceDir = filepath.Join(pl.WorkingDir, filepath.Base(srcResPath))
	logger.Debugf(`Copying resource files: %s`, srcResPath)
	if e := utils.CopyDir(srcResPath, pl.metadata.ResourceDir); e != nil {
		return e
	}

	// Variables
	pl.metadata.Variables = devenv.Variables(pl.Profile)
	pl.metadata.Variables = append(pl.metadata.Variables,
		devenv.Variable{Name: devenv.VarLocalDataPath, Value: pl.Profile.LocalDataDir()},
		devenv.Variable{Name: devenv.VarProjectResource, Value: filepath.Base(srcResPath)},
	)

	// prepare a docker client (this client is not for docker compose)
	pl.dockerClient, e = dockerclient.NewClientWithOpts(dockerclient.WithAPIVersionNegotiation())
	if e != nil {
		return fmt.Errorf("docker client not available: %v", e)
	}

	// docker version
	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	if pl.metadata.DockerVersion, e = pl.dockerClient.ServerVersion(ctx); e != nil {
		return e
	}
	return nil
}

func (pl *DockerComposePlanner) Plan(action Action) (ExecutionPlan, error) {
	if e := pl.Prepare(); e != nil {
		return nil, e
	}
	var execs []Executable
	var e error
	switch action {
	case ActionStart:
		execs, e = pl.startPlan()
	case ActionStop:
		execs, e = pl.stopPlan()
	case ActionRestart:
		execs, e = pl.restartPlan()
	default:
		e = ErrPlanNotAvailable
	}
	if e != nil {
		return nil, e
	}

	return closerPlan{
		ExecutionPlan: NewExecutionPlan(pl.metadata, execs...),
		closeFunc: func() error {
			return pl.dockerClient.Close()
		},
	}, nil
}

func (pl *DockerComposePlanner) startPlan() ([]Executable, error) {
	plan := make([]Executable, 0, 5)
	// step 1 pre-start hooks
	pre, e := pl.hooksPlan(devenv.PhasePreStart, lanaiutils.NewGenericSet(devenv.TypeScript))
	if e != nil {
		return nil, e
	}
	plan = append(plan, pre...)

	// step 2 create data folders if not exist
	dv, e := pl.dataVolumesPlan()
	if e != nil {
		return nil, e
	}
	plan = append(plan, dv...)

	// step 3 docker compose start
	buildCmd := fmt.Sprintf(`docker compose -f "%s" -p "%s" build`, pl.metadata.ComposePath, pl.Profile.Name)
	upCmd := fmt.Sprintf(`docker compose -f "%s" -p "%s" up -d --force-recreate --remove-orphans`, pl.metadata.ComposePath, pl.Profile.Name)
	plan = append(plan, &ShellExecutable{
		Cmds: []string{buildCmd, upCmd},
		WD:   pl.WorkingDir,
		Env:  pl.shellVars(pl.metadata.Variables),
		Desc: "start services",
	})

	// step 4 post-start hooks
	post, e := pl.hooksPlan(devenv.PhasePostStart, lanaiutils.NewGenericSet(devenv.TypeScript, devenv.TypeContainer))
	if e != nil {
		return nil, e
	}
	plan = append(plan, post...)
	return plan, nil
}

func (pl *DockerComposePlanner) stopPlan() ([]Executable, error) {
	return nil, nil
}

func (pl *DockerComposePlanner) restartPlan() ([]Executable, error) {
	return nil, nil
}

func (pl *DockerComposePlanner) hooksPlan(phase devenv.HookPhase, types lanaiutils.GenericSet[devenv.HookType]) ([]Executable, error) {
	hooks, _ := pl.Profile.Hooks[phase]
	execs := make([]Executable, 0, len(hooks))
	vars := pl.shellVars(pl.metadata.Variables)
	var hasContainerHooks bool
	for i := range hooks {
		if !types.Has(hooks[i].Type) {
			return nil, fmt.Errorf(`hooks at phase [%v] only support %v types`, phase, types)
		}

		switch hooks[i].Type {
		case devenv.TypeScript:
			exec, e := pl.shellHookExec(hooks[i], vars)
			if e != nil {
				return nil, e
			}
			execs = append(execs, exec)
		case devenv.TypeContainer:
			hasContainerHooks = true
		default:
			return nil, fmt.Errorf(`unsupported hook type [%v]`, hooks[i].Type)
		}
	}
	// All container hooks are grouped into single executable
	if hasContainerHooks {
		exec, e := pl.containerHookExec(hooks, phase)
		if e != nil {
			return nil, e
		}
		// Note: in post-start, we wait for container at the end. in pre-stop, we wait for container at beginning
		if phase == devenv.PhasePostStop {
			execs = append(execs, exec)
		} else {
			execs = append([]Executable{exec}, execs...)
		}

	}
	return execs, nil
}

func (pl *DockerComposePlanner) dataVolumesPlan() ([]Executable, error) {
	root := pl.Profile.LocalDataDir()
	paths := make([]string, 0, len(pl.Profile.Services)*2)
	for _, s := range pl.Profile.Services {
		for _, mount := range s.Mounts {
			paths = append(paths, filepath.Join(root, mount))
		}
	}
	return []Executable{
		&MkdirExecutable{
			Paths: paths,
			Desc:  "create directories",
		},
	}, nil
}

func (pl *DockerComposePlanner) shellVars(src []devenv.Variable) []string {
	return utils.ConvertSlice[devenv.Variable, string](src, func(v devenv.Variable) string {
		return v.String()
	})
}

func (pl *DockerComposePlanner) shellHookExec(hook devenv.Hook, vars []string) (Executable, error) {
	// note: we assume the value of hook is a script filename in resource directory
	str, ok := hook.Value.(string)
	if !ok {
		return nil, fmt.Errorf(`expected hook value to be string, but got %v`, hook.Value)
	}
	searchPaths := []string{
		filepath.Join(pl.metadata.ResourceDir, string(hook.Phase), str),
		filepath.Join(pl.metadata.ResourceDir, string(hook.Phase)+"-"+str),
	}
	var cmd string
	for _, path := range searchPaths {
		if stat, e := os.Stat(path); e == nil && stat.IsDir() {
			cmd = path
			break
		}
	}
	if len(cmd) == 0 {
		return nil, fmt.Errorf(`hook script [%s] not found in %s`, str, filepath.Join(pl.metadata.ResourceDir, string(hook.Phase), str))
	}
	return &ShellExecutable{
		Cmds: []string{cmd},
		WD:   pl.WorkingDir,
		Env:  vars,
		Desc: fmt.Sprintf(`%v shell`, hook.Phase),
	}, nil
}

func (pl *DockerComposePlanner) containerHookExec(hooks []devenv.Hook, phase devenv.HookPhase) (Executable, error) {
	containers := make([]string, 0, len(hooks))
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
	if phase == devenv.PhasePreStop && len(containers) != 0 {
		// TODO start post container manually
		return nil, fmt.Errorf(`container hook in pre-stop phase is not currently supported`)
	}
	// TODO support dynamic timeout
	exec := NewContainerMonitorExecutable(pl.dockerClient, func(exec *ContainerMonitorExecutable) {
		exec.Names = containers
		exec.Desc = fmt.Sprintf(`%v containers`, phase)
		exec.Resolver = ComposeContainerResolver(pl.Profile.Name)
	})
	return exec.WithTimeout(30 * time.Second), nil
}
