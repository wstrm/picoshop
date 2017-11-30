package view

import (
	"context"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime"
)

const (
	templateExtension = ".tmpl.html"
)

var (
	_, binpath, _, _ = runtime.Caller(0)
	basepath         = filepath.Dir(binpath)
	funcs            = template.FuncMap{
		"context": contextSignature,
	}
	templates = template.Must(template.New("picoshop").
			Funcs(funcs).
			ParseGlob(basepath + "/*" + templateExtension))
)

type Page struct {
	Title string
	User  string
	Data  interface{}
}

func contextSignature(_ interface{}) interface{} {
	return nil
}

func Render(ctx context.Context, writer http.ResponseWriter, tpl string, page Page) {
	funcs := template.FuncMap{
		"context": ctx.Value,
	}

	err := templates.Funcs(funcs).ExecuteTemplate(writer, tpl+".tmpl.html", page)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
