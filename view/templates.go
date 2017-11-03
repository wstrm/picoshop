package view

import (
	"html/template"
	"net/http"
)

const (
	templateExtension = ".tmpl.html"
)

var templates = template.Must(template.ParseGlob("./view/*"+templateExtension))

func Render(writer http.ResponseWriter, template string, data interface{}) {
	err := templates.ExecuteTemplate(writer, template+".tmpl.html", data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}
