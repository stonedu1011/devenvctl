package rootcmd

import (
	"github.com/cisco-open/go-lanai/cmd/lanai-cli/cmdutils"
	"github.com/cisco-open/go-lanai/pkg/log"
	"github.com/cisco-open/go-lanai/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	logTemplate = `{{pad -25 .time}} {{lvl 0 . | printf "[%s]" | pad 23}}: {{.msg}}`
	logProps    = log.Properties{
		Levels: map[string]log.LoggingLevel{
			"default": log.LevelInfo,
		},
		Loggers: map[string]*log.LoggerProperties{
			"console": {
				Type:     log.TypeConsole,
				Format:   log.FormatText,
				Template: logTemplate,
				FixedKeys: utils.CommaSeparatedSlice{
					log.LogKeyName, log.LogKeyMessage, log.LogKeyTimestamp,
					log.LogKeyCaller, log.LogKeyLevel, log.LogKeyContext,
				},
			},
		},
		Mappings: map[string]string{},
	}
)

func init() {
	if e := log.UpdateLoggingConfiguration(&logProps); e != nil {
		panic(e)
	}
}

func New(name, version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                name,
		Version:            version,
		Short:              "A development environment management CLI tool",
		Long:               "This CLI tool allows developers to quickly build, start, stop and switch development environments using Docker",
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		PersistentPreRunE: cmdutils.MergeRunE(
			UpdateLogLevelRunE(),
			cmdutils.EnsureDir(&GlobalArgs.TmpDir, GlobalArgs.WorkingDir, true, "temporary directory"),
			PrintHeaderRunE(),
			SearchProfilesRunE(),
		),
	}
	cmdutils.PersistentFlags(cmd, &GlobalArgs)
	return cmd
}
