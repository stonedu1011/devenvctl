package plan

import (
	"fmt"
	lanaiutils "github.com/cisco-open/go-lanai/pkg/utils"
	"github.com/stonedu1011/devenvctl/pkg/devenv"
	"github.com/stonedu1011/devenvctl/pkg/utils"
	"github.com/stonedu1011/devenvctl/pkg/utils/tmplutils"
	"os"
	"path/filepath"
)

const defaultComposeFile = `docker-compose.yml`

func NewDockerComposePlanner(p *devenv.Profile, wd string) *DockerComposePlanner {
	plan := DockerComposePlanner{
		Profile:    p,
		Variables:  devenv.Variables(p),
		WorkingDir: utils.AbsPath(wd, p.FS),
	}
	return &plan
}

type DockerComposePlanner struct {
	Profile   *devenv.Profile
	Variables []devenv.Variable
	// WorkingDir the working directory. Usually is the temporary dir configured by rootcmd.GlobalArgs
	WorkingDir  string
	resDir      string
	composePath string
}

func (pl *DockerComposePlanner) Prepare() (err error) {
	logger.Infof(`Using Docker Compose`)
	defer func() {
		if err == nil {
			logger.Infof(`Working directory is ready: %s`, pl.WorkingDir)
		}
	}()

	// load compose template
	tmplPath := pl.Profile.ComposePath()
	logger.Debugf(`Loading [%s]`, tmplPath)
	tmpl, e := tmplutils.NewTemplate().ParseFS(pl.Profile.FS, tmplPath)
	if e != nil {
		return fmt.Errorf("unable to process docker compose template [%s]: %v", utils.AbsPath(tmplPath, pl.Profile.FS), e)
	}

	// create a docker-compose.yml
	pl.composePath = filepath.Join(pl.WorkingDir, defaultComposeFile)
	logger.Debugf(`Finalizing docker compose file: %s`, pl.composePath)
	switch fi, e := os.Stat(pl.WorkingDir); {
	case e != nil:
		return fmt.Errorf(`unable to access directory [%s]: %v`, pl.WorkingDir, e)
	case !fi.IsDir():
		return fmt.Errorf(`unable to access directory [%s]: not a directory`, pl.WorkingDir)
	}
	composeF, e := os.OpenFile(pl.composePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if e != nil {
		return fmt.Errorf(`unable to generate docker compose [%s]: %v`, pl.composePath, e)
	}
	defer func() { _ = composeF.Close() }()

	// generate docker-compose.yml
	if e := tmpl.ExecuteTemplate(composeF, filepath.Base(tmplPath), pl); e != nil {
		return fmt.Errorf(`unable to generate docker compose [%s]: %v`, pl.composePath, e)
	}

	// copy resources
	srcResPath := pl.Profile.ResourceDir()
	pl.resDir = filepath.Join(pl.WorkingDir, filepath.Base(srcResPath))
	logger.Debugf(`Copying resource files: %s`, srcResPath)
	return utils.CopyDir(srcResPath, pl.resDir)
}

func (pl *DockerComposePlanner) Plan(action Action) ([]Executable, error) {
	if e := pl.Prepare(); e != nil {
		return nil, e
	}
	switch action {
	case ActionStart:
		return pl.startPlan()
	case ActionStop:
		return pl.stopPlan()
	case ActionRestart:
		return pl.restartPlan()
	default:
		return nil, ErrPlanNotAvailable
	}
}

func (pl *DockerComposePlanner) startPlan() ([]Executable, error) {
	plan := make([]Executable, 0, 5)
	// step 1 pre-start hooks
	pre, e := pl.hooksPlan(devenv.PhasePreStart, lanaiutils.NewGenericSet(devenv.TypeScript))
	if e != nil {
		return nil, e
	}
	plan = append(plan, pre...)
	// TODO
	// step 2 do start
	// step 3 post-start hooks
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
	hooks := pl.Profile.Hooks[phase]
	execs := make([]Executable, 0, len(hooks))
	shellExecTmpl := &ShellExecutable{
		WD:   pl.WorkingDir,
		Env:  pl.vars(pl.Variables),
		Desc: fmt.Sprintf(`%v shell`, phase),
	}
	for _, h := range hooks {
		if !types.Has(h.Type) {
			return nil, fmt.Errorf(`hooks at phase [%v] only support %v types`, phase, types)
		}

		var exec Executable
		var e error
		switch h.Type {
		case devenv.TypeScript:
			exec, e = pl.shellHook(shellExecTmpl, h)
		case devenv.TypeContainer:
			exec, e = pl.containerHook(h)
		default:
			return nil, fmt.Errorf(`unsupported hook type [%v]`, h.Type)
		}

		switch {
		case e != nil:
			return nil, e
		case exec != nil:
			execs = append(execs, exec)
		}
	}
	return execs, nil
}

func (pl *DockerComposePlanner) vars(src []devenv.Variable) []string {
	return utils.ConvertSlice[devenv.Variable, string](src, func(v devenv.Variable) string {
		return v.String()
	})
}

func (pl *DockerComposePlanner) shellHook(execTmpl *ShellExecutable, hook devenv.Hook) (Executable, error) {
	// note: we assume the value of hook is a script filename in resource directory
	str, ok := hook.Value.(string)
	if !ok {
		return nil, fmt.Errorf(`expected hook value to be string, but got %v`, hook.Value)
	}
	searchPaths := []string{
		filepath.Join(pl.resDir, string(hook.Phase), str),
		filepath.Join(pl.resDir, string(hook.Phase)+"-"+str),
	}
	var cmd string
	for _, path := range searchPaths {
		if stat, e := os.Stat(path); e == nil && stat.IsDir() {
			cmd = path
			break
		}
	}
	if len(cmd) == 0 {
		return nil, fmt.Errorf(`hook script [%s] not found in %s`, str, filepath.Join(pl.resDir, string(hook.Phase), str))
	}
	return execTmpl.WithCommands(cmd), nil
}

func (pl *DockerComposePlanner) containerHook(hook devenv.Hook) (Executable, error) {
	// TODO implement this
	return &ShellExecutable{
		Cmds: []string{"echo TODO"},
		WD:   pl.WorkingDir,
		Env:  pl.vars(pl.Variables),
		Desc: fmt.Sprintf(`%v container`, hook.Phase),
	}, nil
}
