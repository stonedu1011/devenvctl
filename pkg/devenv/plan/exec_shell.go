package plan

import (
	"context"
	"fmt"
	"github.com/cisco-open/go-lanai/cmd/lanai-cli/cmdutils"
	"github.com/stonedu1011/devenvctl/pkg/devenv"
	"os"
	"strings"
)

func NewShellVars(src devenv.Variables) []string {
	vars := make([]string, src.Len())
	for i, k := range src.Keys() {
		vars[i] = src.Get(k).String()
	}
	return vars
}


type ShellExecutable struct {
	Cmds []string
	WD   string
	Env  []string
	Desc string
}

func (exec ShellExecutable) Exec(ctx context.Context, opts ExecOption) error {
	if len(exec.Cmds) == 0 {
		return nil
	}
	if opts.DryRun {
		fmt.Printf("- %v\n", exec)
		return nil
	}
	rc, e := cmdutils.RunShellCommands(ctx,
		cmdutils.ShellShowCmd(opts.Verbose),
		cmdutils.ShellDir(exec.WD),
		cmdutils.ShellEnv(exec.Env...),
		cmdutils.ShellCmd(exec.Cmds...),
	)
	switch {
	case e != nil:
		return e
	case rc != 0:
		return fmt.Errorf("shell exited with non-zero code [%d]", rc)
	}
	return nil
}

func (exec ShellExecutable) String() string {
	if len(exec.Desc) == 0 {
		exec.Desc = "shell"
	}
	switch {
	case len(exec.Cmds) == 1:
		return fmt.Sprintf("%s: %s", exec.Desc, exec.Cmds[0])
	case len(exec.Cmds) == 0:
		return "NONE"
	default:
		return fmt.Sprintf("%s: \n    %s", exec.Desc, strings.Join(exec.Cmds, "; \\\n    "))
	}
}

func (exec ShellExecutable) WithCommands(cmds ...string) *ShellExecutable {
	cpy := exec
	cpy.Cmds = cmds
	return &cpy
}

type MkdirExecutable struct {
	Paths         []string
	Desc          string
}

func (exec MkdirExecutable) Exec(ctx context.Context, opts ExecOption) error {
	if len(exec.Paths) == 0 {
		return nil
	}
	if opts.DryRun {
		fmt.Printf("- %v\n", exec)
		return nil
	}
	for _, p := range exec.Paths {
		if opts.Verbose {
			logger.WithContext(ctx).Debugf(`creating direcotry: %s`, p)
		}
		e := os.MkdirAll(p, 0755)
		if e != nil {
			return fmt.Errorf(`unable to create directory [%s]: %v`, p, e)
		}
	}
	return nil
}

func (exec MkdirExecutable) String() string {
	if len(exec.Desc) == 0 {
		exec.Desc = "mkdir"
	}
	switch {
	case len(exec.Paths) == 1:
		return fmt.Sprintf("%s: %s", exec.Desc, exec.Paths[0])
	case len(exec.Paths) == 0:
		return "NONE"
	default:
		return fmt.Sprintf("%s: \n    %s", exec.Desc, strings.Join(exec.Paths, "\n    "))
	}
}

// ComposeShellExecutable ShellExecutable variant with special dry-run strategy for docker compose CLI
type ComposeShellExecutable struct {
	Args []string
	WD   string
	Env  []string
	Desc string
}

func (exec ComposeShellExecutable) Exec(ctx context.Context, opts ExecOption) error {
	if len(exec.Args) == 0 {
		return nil
	}
	cmd := "docker compose " + strings.Join(exec.Args, " ")
	if opts.DryRun {
		cmd = cmd + " --dry-run"
		fmt.Printf("- %v\n", exec)
		opts.Verbose = false
	}

	rc, e := cmdutils.RunShellCommands(ctx,
		cmdutils.ShellShowCmd(opts.Verbose),
		cmdutils.ShellDir(exec.WD),
		cmdutils.ShellEnv(exec.Env...),
		cmdutils.ShellCmd(cmd),
	)
	switch {
	case e != nil:
		return e
	case rc != 0:
		return fmt.Errorf("shell exited with non-zero code [%d]", rc)
	}
	return nil
}

func (exec ComposeShellExecutable) String() string {
	if len(exec.Desc) == 0 {
		exec.Desc = "docker compose"
	}
	switch {
	case len(exec.Args) == 0:
		return "no-op"
	default:
		return fmt.Sprintf("%s: docker compose %s", exec.Desc, strings.Join(exec.Args, " "))
	}
}