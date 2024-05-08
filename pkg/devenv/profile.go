package devenv

import (
	"fmt"
	"github.com/stonedu1011/devenvctl/pkg/utils"
	"github.com/stonedu1011/devenvctl/pkg/utils/tmplutils"
	"io/fs"
)

type Profiles map[string]*ProfileMetadata

type ProfileMetadata struct {
	FS          fs.FS  `json:"-"`
	Name        string `json:"-"`
	Path        string `json:"-"`
	DisplayPath string `json:"-"`
}

type Profile struct {
	ProfileMetadata
	DisplayName string
	Services    []Service
	Hooks       Hooks
}

func (p Profile) ResourceDir() string {
	return tmplutils.MustSprint(ResourceDirTemplate, p)
}

func (p Profile) ComposePath() string {
	return tmplutils.MustSprint(ComposeTemplate, p)
}

func FindProfiles(fsys fs.FS, searchPaths ...string) (Profiles, error) {
	fsysName := fmt.Sprintf(`%v`, fsys)
	profiles := Profiles{}
	for _, searchPath := range searchPaths {
		logger.Debugf(`Searching [%s] ...`, utils.AbsPath(searchPath, fsysName))
		e := fs.WalkDir(fsys, searchPath, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}
			displayPath := utils.AbsPath(path, fsysName)
			var matched bool
			defer func() {
				if matched {
					logger.Debugf(`Matched %s`, displayPath)
				} else {
					logger.Debugf(`Ignored %s`, displayPath)
				}
			}()

			fn := d.Name()
			for _, regex := range RegexProfile {
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
