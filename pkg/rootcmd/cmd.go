package rootcmd

import (
	"github.com/cisco-open/go-lanai/cmd/lanai-cli/cmdutils"
	"github.com/cisco-open/go-lanai/pkg/log"
	"github.com/cisco-open/go-lanai/pkg/utils"
	"github.com/spf13/cobra"
)

const (
	description = `
This CLI tool allows developers to quickly build, start, stop and switch development environments using Docker Compose.

Environment Profiles:

Environment profiles are defined in YAML with filename "devenv-<profile-name>.yml". 
This tool will search profile definitions in following order:
    - "~/.devenv" or "$HOME/.devenv"
    - working directory specified via "--workspace", default to current directory
    - $DEV_ENV_PATH
    - any additional search path defined via "--search-paths"
`
)

var (
	logTemplate        = `{{pad -25 .time}} [{{lvl 4 .}}]: {{.msg}}`
	logVerboseTemplate = `{{pad -25 .time}} [{{lvl 5 .}}]: {{.msg}}`
)

func init() {
	MustUpdateLoggingConfiguration(NewLogConfig(log.LevelInfo, logTemplate))
	cobra.OnInitialize(func() {
		if GlobalArgs.Verbose {
			MustUpdateLoggingConfiguration(NewLogConfig(log.LevelDebug, logVerboseTemplate))
		}
	})
}

func New(name string) *cobra.Command {
	cmd := &cobra.Command{
		Use:                name,
		Short:              "A development environment management CLI tool",
		Long:               description,
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		PersistentPreRunE: cmdutils.MergeRunE(
			cmdutils.EnsureDir(&GlobalArgs.TmpDir, GlobalArgs.WorkingDir, true, "temporary directory"),
			PrintHeaderRunE(),
			SearchProfilesRunE(),
		),
	}
	cmdutils.PersistentFlags(cmd, &GlobalArgs)
	return cmd
}

func MustUpdateLoggingConfiguration(props *log.Properties) {
	if e := log.UpdateLoggingConfiguration(props); e != nil {
		panic(e)
	}
}

func NewLogConfig(lvl log.LoggingLevel, tmplText string) *log.Properties {
	return &log.Properties{
		Levels: map[string]log.LoggingLevel{
			"default": lvl,
		},
		Loggers: map[string]*log.LoggerProperties{
			"console": {
				Type:     log.TypeConsole,
				Format:   log.FormatText,
				Template: tmplText,
				FixedKeys: utils.CommaSeparatedSlice{
					log.LogKeyName, log.LogKeyMessage, log.LogKeyTimestamp,
					log.LogKeyCaller, log.LogKeyLevel, log.LogKeyContext,
				},
			},
		},
		Mappings: map[string]string{},
	}
}