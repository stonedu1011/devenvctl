package debug

import (
	"fmt"
	"github.com/cisco-open/go-lanai/cmd/lanai-cli/cmdutils"
	"github.com/spf13/cobra"
)

const (
	CommandName = "debug"
)

var (
	Cmd = &cobra.Command{
		Use:                fmt.Sprintf(`%s ...`, CommandName),
		Short:              "Internal use only",
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		//Args:               cobra.MinimumNArgs(1),
		//PreRunE:            rootcmd.LoadProfileRunE(),
		RunE:               Run,
		Hidden:             true,
	}
	Args = Arguments{}
)

type Arguments struct {
	DryRun bool `flag:"dry-run" desc:"print out commands instead of run them"`
}

func init() {
	cmdutils.PersistentFlags(Cmd, &Args)
}

func Run(cmd *cobra.Command, args []string) error {

	return nil
}
