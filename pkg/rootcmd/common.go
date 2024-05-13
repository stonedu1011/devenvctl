package rootcmd

import (
	"errors"
	"fmt"
	"github.com/cisco-open/go-lanai/cmd/lanai-cli/cmdutils"
	"github.com/cisco-open/go-lanai/pkg/log"
	"github.com/spf13/cobra"
	"github.com/stonedu1011/devenvctl/pkg/devenv"
	"github.com/stonedu1011/devenvctl/pkg/tmpls"
	"github.com/stonedu1011/devenvctl/pkg/utils/tmplutils"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var logger = log.New("CLI")

var searchOnce sync.Once

var (
	Profiles      devenv.Profiles
	LoadedProfile *devenv.Profile
)

var (
	GlobalArgs = Global{
		WorkingDir: DefaultWorkingDir(),
		TmpDir:     DefaultTemporaryDir(),
	}
)

type Global struct {
	WorkingDir  string   `flag:"workspace,w" desc:"working directory containing profile definitions"`
	TmpDir      string   `flag:"tmp-dir" desc:"temporary directory."`
	Verbose     bool     `flag:"verbose,v" desc:"show debug information"`
	SearchPaths []string `flag:"search-paths,s" desc:"additional paths to search for profiles definitions"`
}

func DefaultWorkingDir() string {
	path, e := os.Getwd()
	if e != nil {
		return "."
	}
	return path
}

func DefaultTemporaryDir() string {
	const tmpDir = `.tmp`
	path, e := os.Getwd()
	if e != nil {
		return tmpDir
	}
	return filepath.Join(path, tmpDir)
}

func RequireProfileArgs() cobra.PositionalArgs {
	return func(_ *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("missing environment's profile name")
		}
		if _, e := SearchProfiles(); e != nil {
			return e
		}
		if _, ok := Profiles[args[0]]; !ok {
			return fmt.Errorf(`unknown profile [%s]`, args[0])
		}
		return nil
	}
}

func PrintHeaderRunE() cmdutils.RunE {
	return func(cmd *cobra.Command, args []string) error {
		tmplData := map[string]interface{}{
			"Cmd":    cmd,
			"Args":   strings.Join(args, " "),
			"Global": GlobalArgs,
		}
		return tmplutils.Print(tmpls.OutputTemplate.Lookup("header.tmpl"), tmplData)
	}
}

func SearchProfilesRunE() cmdutils.RunE {
	return func(cmd *cobra.Command, args []string) (err error) {
		_, err = SearchProfiles()
		return
	}
}

// LoadProfileRunE common RunE for any command that requires profile as argument
func LoadProfileRunE() cmdutils.RunE {
	return func(cmd *cobra.Command, args []string) error {
		// Arguments should be verified at this moment
		pName := args[0]
		profiles, e := SearchProfiles()
		if e != nil {
			return e
		}
		pMeta := profiles[pName]
		LoadedProfile, e = devenv.LoadProfile(pMeta)
		if e != nil {
			return e
		}
		if e := tmplutils.Print(tmpls.OutputTemplate.Lookup("profile.tmpl"), LoadedProfile); e != nil {
			return e
		}

		if e := tmplutils.Print(tmpls.OutputTemplate.Lookup("mounts.tmpl"), LoadedProfile); e != nil {
			return e
		}
		return nil
	}
}
