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
	TemplateResourceDir      = tmplutils.MustParse(`{{.Dir}}/res-{{.Name}}`)
	TemplateComposePath      = tmplutils.MustParse(`{{.Dir}}/docker-compose-{{.Name}}.yml`)
	TemplateLocalDataDir     = tmplutils.MustParse(`/usr/local/var/dev/{{.Name}}`)
	TemplateBuildArgName     = tmplutils.MustParse(`build_args_{{.}}`)
	TemplateServiceImage     = tmplutils.MustParse(`{{.Name}}_image`)
	TemplateServiceContainer = tmplutils.MustParse(`{{.Name}}_container_name`)
)

const (
	VarProjectName     = `PROJECT_NAME`
	VarProjectResource = `RESOURCE_DIR`
	VarLocalDataPath   = `CONTAINER_DATA_PATH`
)
