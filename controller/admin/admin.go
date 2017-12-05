package admin

import (
	"log"
	"net/http"

	"github.com/willeponken/picoshop/view"
)

type rootHandler struct {
	http.Handler
}

type registerHandler struct {
	http.Handler
}

type articleHandler struct {
	http.Handler
}

type rootData struct {
	Error string
	// Admin     registerData
	// Warehouse registerData
	// Customer  registerData
}

func (a *rootHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	//ctx := request.Context()
	log.Println(request.URL)

	switch request.Method {
	case http.MethodGet:
		view.Render(request.Context(), writer, "admin", view.Page{Title: "Admin - Picoshop", Data: rootData{}})
	}
}

func (a *registerHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		view.Render(request.Context(), writer, "admin", view.Page{Title: "Admin - Picoshop", Data: rootData{}})
	}
}

func (a *articleHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		view.Render(request.Context(), writer, "admin", view.Page{Title: "Admin - Picoshop", Data: rootData{}})
	}
}

func newRootHandler() *rootHandler {
	return &rootHandler{}
}

func newArticleHandler() *articleHandler {
	return &articleHandler{}
}

func newRegisterHandler() *registerHandler {
	return &registerHandler{}
}

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", newRootHandler())
	mux.Handle("/register", newRegisterHandler())
	mux.Handle("/article", newArticleHandler())

	return mux
}
