package devenv

import (
	"github.com/cisco-open/go-lanai/pkg/log"
	"github.com/stonedu1011/devenvctl/pkg/utils/tmplutils"
	"regexp"
)

var logger = log.New("CLI")

var (
	RegexProfile = []*regexp.Regexp{
		regexp.MustCompile(`devenv-(?P<profile>[a-zA-Z][\w-_]+)\.yml`),
		regexp.MustCompile(`devenv-(?P<profile>[a-zA-Z][\w-_]+)\.yaml`),
	}
)

var (
	ResourceDirTemplate = tmplutils.MustParse(`res-{{.Name}}`)
	ComposeTemplate = tmplutils.MustParse(`docker-compose-{{.Name}}.yml`)
	BuildArgNameTemplate = tmplutils.MustParse(`build_args_{{.}}`)
)
