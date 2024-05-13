package devenv

import (
	"github.com/cisco-open/go-lanai/pkg/log"
	"github.com/stonedu1011/devenvctl/pkg/utils/tmplutils"
)

var logger = log.New("CLI")

var (
	TemplateBuildArgName     = tmplutils.MustParse(`build_args_{{.}}`)
	TemplateServiceImage     = tmplutils.MustParse(`{{.Name}}_image`)
	TemplateServiceContainer = tmplutils.MustParse(`{{.Name}}_container_name`)
)

const (
	VarProjectName     = `PROJECT_NAME`
	VarProjectResource = `RESOURCE_DIR`
	VarLocalDataPath   = `CONTAINER_DATA_PATH`
)
