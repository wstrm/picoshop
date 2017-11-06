package register

import (
	"net/http"

	"github.com/willeponken/picoshop/view"
)

type handler struct {
	http.Handler
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {

	case http.MethodGet: // Serve register view
		view.Render(writer, "register", view.Page{
			Title: "Picoshop",
		})

	case http.MethodPost: // Retreive user registration
		http.Error(writer, "Not Implemented", http.StatusNotImplemented)
	}
}

func New() *handler {
	return &handler{}
}
