package controller

import (
	"net/http"

	"github.com/willeponken/picoshop/view"
)

type homeHandler struct {
	http.Handler
}

func (h *homeHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}

	view.Render(writer, "home", view.Page{Title: "Picoshop"})
}

func newHomeHandler() *homeHandler {
	return &homeHandler{}
}
