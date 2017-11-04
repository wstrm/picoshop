package home

import (
	"net/http"

	"github.com/willeponken/picoshop/view"
)

type handler struct {
	http.Handler
}

type page struct {
	Title string
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}

	view.Render(writer, "home", &page{Title: "Picoshop"})
}

func New() *handler {
	return &handler{}
}
