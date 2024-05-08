package info

import (
	"embed"
	"fmt"
	"github.com/cisco-open/go-lanai/cmd/lanai-cli/cmdutils"
	"github.com/spf13/cobra"
	"github.com/stonedu1011/devenvctl/pkg/devenv"
	"github.com/stonedu1011/devenvctl/pkg/rootcmd"
	"github.com/stonedu1011/devenvctl/pkg/utils/tmplutils"
)

const (
	RootName = "info"
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
	//Metadata string `flag:"module-metadata,m" desc:"metadata yaml for the module"`
}

func init() {
	cmdutils.PersistentFlags(Cmd, &Args)
}

//go:embed *.tmpl
var outputTmplFS embed.FS

func Run(_ *cobra.Command, _ []string) error {
	if !rootcmd.GlobalArgs.Verbose {
		return nil
	}

	vars := devenv.VariablesBuildArg(rootcmd.LoadedProfile)
	if e := tmplutils.PrintFS(outputTmplFS, "build_args.tmpl", vars); e != nil {
		return e
	}

	if e := tmplutils.PrintFS(outputTmplFS, "hooks.tmpl", rootcmd.LoadedProfile); e != nil {
		return e
	}
	return nil
}
