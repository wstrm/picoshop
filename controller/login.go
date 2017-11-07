package controller

import (
	"net/http"

	"github.com/willeponken/picoshop/view"
)

type loginHandler struct {
	http.Handler
}

func (l *loginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		view.Render(writer, "login", view.Page{Title: "Picoshop"})

	case http.MethodPost:
		http.Error(writer, "", http.StatusNotImplemented)

	default:
		http.Error(writer, "", http.StatusMethodNotAllowed)
	}
}

func newLoginHandler() *loginHandler {
	return &loginHandler{}
}
