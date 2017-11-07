package controller

import (
	"net/http"

	"github.com/willeponken/picoshop/view"
)

type registerHandler struct {
	http.Handler
}

func (r *registerHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {

	case http.MethodGet: // Serve register view
		view.Render(writer, "register", view.Page{
			Title: "Picoshop",
		})

	case http.MethodPost: // Retreive user registration
		http.Error(writer, "Not Implemented", http.StatusNotImplemented)
	}
}

func newRegisterHandler() *registerHandler {
	return &registerHandler{}
}
