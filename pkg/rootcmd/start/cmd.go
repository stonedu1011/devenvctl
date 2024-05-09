package start

import (
	"fmt"
	"github.com/cisco-open/go-lanai/cmd/lanai-cli/cmdutils"
	"github.com/spf13/cobra"
	"github.com/stonedu1011/devenvctl/pkg/devenv/plan"
	"github.com/stonedu1011/devenvctl/pkg/rootcmd"
	"github.com/stonedu1011/devenvctl/pkg/utils"
	"github.com/stonedu1011/devenvctl/pkg/utils/tmplutils"
)

const (
	RootName = "start"
)

var (
	Cmd = &cobra.Command{
		Use:                fmt.Sprintf(`%s <profile>`, RootName),
		Short:              "Show information of specified profile",
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Args:               rootcmd.RequireProfileArgs(),
		PreRunE:            rootcmd.LoadProfileRunE(),
		RunE:               Run,
	}
	Args = Arguments{}
)

type Arguments struct {
	DryRun bool `flag:"dry-run" desc:"print out commands instead of run them"`
}

func init() {
	cmdutils.PersistentFlags(Cmd, &Args)
}

func Run(cmd *cobra.Command, _ []string) error {
	tmpDir := utils.AbsPath(rootcmd.GlobalArgs.TmpDir, rootcmd.GlobalArgs.WorkingDir)
	planner := plan.NewDockerComposePlanner(rootcmd.LoadedProfile, tmpDir)
	execs, e := planner.Plan(plan.ActionStart)
	if e != nil {
		return e
	}

	if rootcmd.GlobalArgs.Verbose {
		if e := tmplutils.PrintFS(rootcmd.OutputTmplFS, "context.tmpl", planner); e != nil {
			return e
		}

		if e := tmplutils.PrintFS(rootcmd.OutputTmplFS, "hooks.tmpl", rootcmd.LoadedProfile); e != nil {
			return e
		}
	}

	if Args.DryRun {
		e = plan.DryRun(cmd.Context(), execs...)
	} else {
		e = plan.Execute(cmd.Context(), execs...)
	}

	return e
}
