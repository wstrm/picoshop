package controller

import (
	"net/http"

	"github.com/willeponken/picoshop/view"
)

type adminHandler struct {
	http.Handler
}

type adminData struct {
	Error string
}

func (a *adminHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	view.Render(request.Context(), writer, "admin", view.Page{Title: "Admin - Picoshop", Data: adminData{}})
}

func newAdminHandler() *adminHandler {
	return &adminHandler{}
}
