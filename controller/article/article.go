package article

import (
	"net/http"

	"github.com/willeponken/picoshop/view"
)

type handler struct {
	http.Handler
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")

	if id != "" {
		// Get article data using ID from model here
		view.Render(writer, "article", view.Page{Title: "Article - Picoshop", Data: id})
		return
	}

	http.NotFound(writer, request)
}

func New() *handler {
	return &handler{}
}
