package tmplutils

import (
	"bytes"
	"fmt"
	"github.com/stonedu1011/devenvctl/pkg/utils/tmplutils/internal"
	"io/fs"
	"os"
	"text/template"
)

type tmplType interface {
	*template.Template | string | []byte
}

func NewTemplate() *template.Template {
	return template.New("template").
		Option("missingkey=zero").
		Funcs(internal.TmplFuncMap).
		Funcs(internal.TmplColorFuncMap)
}

func MustParse(tmplText string) *template.Template {
	t, e := NewTemplate().Parse(tmplText)
	if e != nil {
		panic(e)
	}
	return t
}

func PrintFS(fsys fs.FS, tmplPath string, data interface{}) error {
	tmpl, e := NewTemplate().ParseFS(fsys, tmplPath)
	if e != nil {
		return e
	}
	return tmpl.ExecuteTemplate(os.Stdout, tmplPath, data)
}

func Print[T tmplType](rawTmpl T, data interface{}) error {
	tmpl, e := parseTemplate(rawTmpl)
	if e != nil {
		return e
	}
	return tmpl.Execute(os.Stdout, data)
}

func MustSprint[T tmplType](rawTmpl T, data interface{}) string {
	s, e := Sprint(rawTmpl, data)
	if e != nil {
		panic(e)
	}
	return s
}

func Sprint[T tmplType](rawTmpl T, data interface{}) (string, error) {
	tmpl, e := parseTemplate(rawTmpl)
	if e != nil {
		return "", e
	}
	var buf bytes.Buffer
	if e := tmpl.Execute(&buf, data); e != nil {
		return "", e
	}
	return buf.String(), nil
}

func parseTemplate[T tmplType](raw T) (t *template.Template, e error) {
	switch v := any(raw).(type) {
	case *template.Template:
		t = v
	case string:
		t, e = NewTemplate().Parse(v)
	case []byte:
		t, e = NewTemplate().Parse(string(v))
	default:
		e = fmt.Errorf(`unsupported type of raw template: %T`, raw)
	}
	return
}
