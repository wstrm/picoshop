package home

import (
	"net/http"

	"github.com/willeponken/picoshop/view"
)

type handler struct {
	http.Handler
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}

	view.Render(writer, "home", view.Page{Title: "Picoshop"})
}

func New() *handler {
	return &handler{}
}
