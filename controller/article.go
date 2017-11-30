package controller

import (
	"net/http"

	"github.com/willeponken/picoshop/view"
)

type articleHandler struct {
	http.Handler
}

func (a *articleHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")

	if id != "" {
		// Get article data using ID from model here
		view.Render(request.Context(), writer, "article", view.Page{Title: "Article - Picoshop", Data: id})
		return
	}

	http.NotFound(writer, request)
}

func newArticleHandler() *articleHandler {
	return &articleHandler{}
}
