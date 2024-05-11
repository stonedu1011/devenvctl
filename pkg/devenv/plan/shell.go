package plan

import (
	"context"
	"fmt"
	"github.com/cisco-open/go-lanai/cmd/lanai-cli/cmdutils"
	"os"
	"strings"
)

type ShellExecutable struct {
	Cmds []string
	WD   string
	Env  []string
	Desc string
}

func (exec ShellExecutable) Exec(ctx context.Context) error {
	if len(exec.Cmds) == 0 {
		return nil
	}
	rc, e := cmdutils.RunShellCommands(ctx,
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
		exec.Desc = "Shell"
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

func (exec MkdirExecutable) Exec(_ context.Context) error {
	if len(exec.Paths) == 0 {
		return nil
	}
	for _, p := range exec.Paths {
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