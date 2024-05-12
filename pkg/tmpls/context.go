package tmpls

import (
	"embed"
	"github.com/stonedu1011/devenvctl/pkg/utils/tmplutils"
)

//go:embed *.tmpl
var OutputTmplFS embed.FS

var OutputTemplate = tmplutils.MustParseGlob(OutputTmplFS, "*.tmpl")