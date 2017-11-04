package register

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
	view.Render(writer, "register", &page{Title: "Picoshop"})
}

func New() *handler {
	return &handler{}
}
