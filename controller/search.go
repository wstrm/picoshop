package controller

import (
	"errors"
	"net/http"

	"github.com/willeponken/picoshop/model"
	"github.com/willeponken/picoshop/view"
)

type searchHandler struct {
	http.Handler
}

type searchData struct {
	Error    string
	Articles []model.Article
}

func renderSearch(writer http.ResponseWriter, code int, data interface{}) {
	writer.WriteHeader(code)

	if data == nil {
		data = searchData{}
	}

	page := view.Page{
		Title: "Search - Picoshop",
		Data:  data,
	}

	view.Render(writer, "search", page)
}

func (a *searchHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query().Get("query")

	articles, err := model.SearchForArticles(query)
	if err != nil {
		renderSearch(writer, http.StatusInternalServerError, searchData{
			Error: errors.New("Something internal went wrong!").Error(),
		})
		return
	}

	renderSearch(writer, http.StatusOK, searchData{
		Articles: articles,
	})
	return
}

func newSearchHandler() *searchHandler {
	return &searchHandler{}
}
