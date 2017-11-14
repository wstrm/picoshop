package view

import (
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
	templates        = template.Must(template.ParseGlob(basepath + "/*" + templateExtension))
)

type Page struct {
	Title string
	Data  interface{}
}

func Render(writer http.ResponseWriter, template string, page Page) {
	err := templates.ExecuteTemplate(writer, template+".tmpl.html", page)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
