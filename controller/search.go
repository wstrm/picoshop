package controller

import (
	"net/http"

	"github.com/willeponken/picoshop/view"
	"github.com/willeponken/picoshop/model"
)

type searchHandler struct {
	http.Handler
}

func (a *searchHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	search := request.URL.Query().Get("search")

	result, _ := model.GetSearchResult(search)
			// Get article data using ID from model here
		view.Render(writer, "search", view.Page{Title: "Search - Result", Data: result})
		return
}

func newSearchHandler() *searchHandler {
	return &searchHandler{}
}

