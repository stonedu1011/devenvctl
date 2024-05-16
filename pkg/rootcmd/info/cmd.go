package info

import (
	"fmt"
	"github.com/cisco-open/go-lanai/cmd/lanai-cli/cmdutils"
	"github.com/spf13/cobra"
	"github.com/stonedu1011/devenvctl/pkg/devenv"
	"github.com/stonedu1011/devenvctl/pkg/rootcmd"
	"github.com/stonedu1011/devenvctl/pkg/tmpls"
	"github.com/stonedu1011/devenvctl/pkg/utils/tmplutils"
)

const (
	CommandName = "info"
)

var (
	Cmd = &cobra.Command{
		Use:                fmt.Sprintf(`%s <profile>`, CommandName),
		Short:              "Show information of specified profile",
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Args:               rootcmd.RequireProfileArgs(),
		PreRunE:            rootcmd.LoadProfileRunE(),
		RunE:               Run,
	}
	Args = Arguments{}
)

type Arguments struct {
	//Metadata string `flag:"module-metadata,m" desc:"metadata yaml for the module"`
}

func init() {
	cmdutils.PersistentFlags(Cmd, &Args)
}

func Run(_ *cobra.Command, _ []string) error {
	if !rootcmd.GlobalArgs.Verbose {
		return nil
	}

	vars := devenv.ResolveBuildArgs(rootcmd.LoadedProfile)
	if e := tmplutils.Print(tmpls.OutputTemplate.Lookup("build_args.tmpl"), vars); e != nil {
		return e
	}

	if e := tmplutils.Print(tmpls.OutputTemplate.Lookup("hooks.tmpl"), rootcmd.LoadedProfile); e != nil {
		return e
	}
	return nil
}
