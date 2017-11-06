package user

import (
	"net/http"

	"github.com/willeponken/picoshop/view"
)

type handler struct {
	http.Handler
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		view.Render(writer, "user", view.Page{Title: "Picoshop"})

	case http.MethodPost:
		http.Error(writer, "", http.StatusNotImplemented)

	default:
		http.Error(writer, "", http.StatusMethodNotAllowed)
	}
}

func New() *handler {
	return &handler{}
}
