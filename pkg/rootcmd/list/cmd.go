package list

import (
	"embed"
	"github.com/cisco-open/go-lanai/cmd/lanai-cli/cmdutils"
	"github.com/spf13/cobra"
	"github.com/stonedu1011/devenvctl/pkg/rootcmd"
	"github.com/stonedu1011/devenvctl/pkg/utils/tmplutils"
)

const (
	RootName = "list"
)

var (
	Cmd    = &cobra.Command{
		Use:                RootName,
		Short:              "List available profiles",
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		Args:               cobra.NoArgs,
		RunE:               Run,
	}
	Args = Arguments{}
)

//go:embed output.tmpl
var templateFS embed.FS

type Arguments struct {
	//Metadata string `flag:"module-metadata,m" desc:"metadata yaml for the module"`
}

func init() {
	cmdutils.PersistentFlags(Cmd, &Args)
}

func Run(_ *cobra.Command, _ []string) error {
	return tmplutils.PrintFS(templateFS, "output.tmpl", map[string]interface{}{
		"Profiles": rootcmd.Profiles,
	})
}
