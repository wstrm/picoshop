package controller

import (
	"net/http"

	"github.com/willeponken/picoshop/view"
)

type userHandler struct {
	http.Handler
}

func (u *userHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		view.Render(request.Context(), writer, "user", view.Page{Title: "Picoshop"})

	case http.MethodPost:
		http.Error(writer, "", http.StatusNotImplemented)

	default:
		http.Error(writer, "", http.StatusMethodNotAllowed)
	}
}

func newUserHandler() *userHandler {
	return &userHandler{}
}
