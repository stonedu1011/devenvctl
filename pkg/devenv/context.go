package devenv

import (
	"github.com/cisco-open/go-lanai/pkg/log"
	"github.com/stonedu1011/devenvctl/pkg/utils/tmplutils"
)

var logger = log.New("CLI")

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
