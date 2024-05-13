package devenv

import (
	"fmt"
	"github.com/stonedu1011/devenvctl/pkg/utils"
	"io/fs"
	"path/filepath"
	"regexp"
)

type Profiles map[string]*ProfileMetadata

type ProfileMetadata struct {
	FS           fs.FS  `json:"-"`
	Name         string `json:"-"`
	Path         string `json:"-"`
	Dir          string `json:"-"`
	DisplayPath  string `json:"-"`
	ResourceDir  string `json:"-"`
	ComposePath  string `json:"-"`
	LocalDataDir string `json:"-"`
}

type Profile struct {
	ProfileMetadata
	DisplayName string
	Services    map[string]Service
	Hooks       Hooks
}

func MergeProfiles(src, dest Profiles) Profiles {
	if dest == nil {
		dest = Profiles{}
	}
	for k, v := range src {
		dest[k] = v
	}
	return dest
}

func FindProfiles(fsys fs.FS, dirPath string, nameRegexps ...*regexp.Regexp) (Profiles, error) {
	profiles := Profiles{}
	logger.Debugf(`Searching [%s] ...`, utils.AbsPath(dirPath, fsys))
	e := fs.WalkDir(fsys, dirPath, func(path string, d fs.DirEntry, err error) error {
		displayPath := utils.AbsPath(path, fsys)
		switch {
		case err != nil:
			logger.Debugf(`Ignoring [%s]: %v`, displayPath, err)
			fallthrough
		case d.IsDir():
			return nil
		}
		var matched bool
		defer func() {
			if matched {
				logger.Debugf(`Matched [%s]`, displayPath)
			}
		}()

		fn := d.Name()
		for _, regex := range nameRegexps {
			matchIdx := regex.SubexpIndex("profile")
			if matches := regex.FindStringSubmatch(fn); len(matches) > matchIdx {
				matched = true
				name := matches[matchIdx]
				if existing, ok := profiles[name]; ok {
					return fmt.Errorf(`found multiple definition files of same profile "%s": "%s" and "%s"`, name, existing, path)
				}
				profiles[name] = &ProfileMetadata{
					FS:          fsys,
					Name:        name,
					Path:        path,
					Dir:         filepath.Dir(path),
					DisplayPath: displayPath,
				}
				break
			}
		}
		return nil
	})
	if e != nil {
		return nil, fmt.Errorf("failed to search profiles: %v", e)
	}

	return profiles, nil
}

func LoadProfile(meta *ProfileMetadata) (*Profile, error) {
	pv1, err := LoadProfileV1(meta)
	if err == nil {
		return pv1.ToProfile(), nil
	}
	// more version format goes here
	return nil, err
}
