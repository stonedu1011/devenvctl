package rootcmd

import (
	"fmt"
	"github.com/stonedu1011/devenvctl/pkg/devenv"
	"github.com/stonedu1011/devenvctl/pkg/utils"
	"io/fs"
	"os"
	"regexp"
)

const (
	EnvSearchPath     = `DEV_ENV_PATH`
	RelWDSearchPath   = `.`
	RelHomeSearchPath = `.devenv`
)

var (
	RegexProfile = []*regexp.Regexp{
		regexp.MustCompile(`devenv-(?P<profile>[a-zA-Z][\w-_]+)\.yml`),
		regexp.MustCompile(`devenv-(?P<profile>[a-zA-Z][\w-_]+)\.yaml`),
	}
)

func SearchProfiles() (loaded devenv.Profiles, err error) {
	searchOnce.Do(func() {
		searchPaths := resolveProfileSources()
		for _, s := range searchPaths {
			loaded, err = devenv.FindProfiles(s.fsys, s.dir, s.regexps...)
			if err != nil {
				return
			}
			Profiles = devenv.MergeProfiles(loaded, Profiles)
		}
	})
	if err == nil && len(Profiles) == 0 {
		return nil, fmt.Errorf(`unable to find profiles`)
	}
	return Profiles, err
}

type profileSource struct {
	fsys    fs.FS
	dir     string
	regexps []*regexp.Regexp
}

func resolveProfileSources() []profileSource {
	srcs := make([]profileSource, 0, 3)
	homeDir, _ := os.UserHomeDir()
	// Home directory
	if len(homeDir) != 0 {
		srcs = append(srcs, profileSource{
			fsys:    os.DirFS(homeDir),
			dir:     RelHomeSearchPath,
			regexps: RegexProfile,
		})
	}

	// working directory
	srcs = append(srcs, profileSource{
		fsys:    os.DirFS(GlobalArgs.WorkingDir),
		dir:     RelWDSearchPath,
		regexps: RegexProfile,
	})

	// From ENV
	if v := os.Getenv(EnvSearchPath); len(v) != 0 {
		srcs = append(srcs, profileSource{
			fsys:    os.DirFS(utils.AbsPath(v, GlobalArgs.WorkingDir)),
			dir:     ".",
			regexps: RegexProfile,
		})
	} else {
		logger.Infof(`$%s is not set`, EnvSearchPath)
	}

	// Additional Paths
	if len(GlobalArgs.SearchPaths) > 0 {
		for _, p := range GlobalArgs.SearchPaths {
			srcs = append(srcs, profileSource{
				fsys:    os.DirFS(utils.AbsPath(p, GlobalArgs.WorkingDir)),
				dir:     ".",
				regexps: RegexProfile,
			})
		}
	}

	return srcs
}
