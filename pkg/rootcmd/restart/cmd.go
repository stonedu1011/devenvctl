package restart

import (
	"fmt"
	"github.com/cisco-open/go-lanai/cmd/lanai-cli/cmdutils"
	"github.com/spf13/cobra"
	"github.com/stonedu1011/devenvctl/pkg/devenv/plan"
	"github.com/stonedu1011/devenvctl/pkg/rootcmd"
	"github.com/stonedu1011/devenvctl/pkg/tmpls"
	"github.com/stonedu1011/devenvctl/pkg/utils"
	"github.com/stonedu1011/devenvctl/pkg/utils/tmplutils"
)

const (
	CommandName = "restart"
)

var (
	Cmd = &cobra.Command{
		Use:                fmt.Sprintf(`%s <profile>`, CommandName),
		Short:              "Stop profile",
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
	p, e := planner.Plan(plan.ActionRestart)
	if e != nil {
		return e
	}

	if rootcmd.GlobalArgs.Verbose {
		if e := tmplutils.Print(tmpls.OutputTemplate.Lookup("docker_plan.tmpl"), p.Metadata()); e != nil {
			return e
		}
	}

	return p.Execute(cmd.Context(), func(opt *plan.ExecOption) {
		opt.DryRun = Args.DryRun
		opt.Verbose = rootcmd.GlobalArgs.Verbose
	})
}
