package home

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

	view.Render(request.Context(), writer, "home", view.Page{Title: "Picoshop"})
}

func NewHandler() *homeHandler {
	return &homeHandler{}
}
