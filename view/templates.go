package view

import (
	"html/template"
	"net/http"
)

const (
	templateExtension = ".tmpl.html"
)

var templates = template.Must(template.ParseGlob("./view/*" + templateExtension))

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
