package rootcmd

import (
	"fmt"
	"github.com/stonedu1011/devenvctl/pkg/devenv"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

const (
	EnvSearchPath = `DEV_ENV_PATH`
	RelWDSearchPath = `devenv`
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
	// From ENV
	if v := os.Getenv(EnvSearchPath); len(v) != 0 {
		absPath := v
		if !filepath.IsAbs(v) {
			absPath = filepath.Join(GlobalArgs.WorkingDir, v)
		}
		srcs = append(srcs, profileSource{
			fsys:    os.DirFS(absPath),
			dir:     ".",
			regexps: RegexProfile,
		})
	} else if len(homeDir) != 0 {
		logger.Warnf(`$%s is not set. Using [%s/%s] and [%s/%s]`, EnvSearchPath, homeDir, RelHomeSearchPath, GlobalArgs.WorkingDir, RelWDSearchPath)
	} else {
		logger.Warnf(`$%s is not set. Using [%s/%s]`, EnvSearchPath, GlobalArgs.WorkingDir, RelWDSearchPath)
	}

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

	return srcs
}
