package search

import (
	"context"
	"errors"
	"net/http"

	"github.com/wstrm/picoshop/model"
	"github.com/wstrm/picoshop/view"
)

type searchHandler struct {
	http.Handler
}

type searchData struct {
	Error    string
	Query    string
	Articles []model.Article
}

func renderSearch(ctx context.Context, writer http.ResponseWriter, code int, data interface{}) {
	writer.WriteHeader(code)

	if data == nil {
		data = searchData{}
	}

	page := view.Page{
		Title: "Search - Picoshop",
		Data:  data,
	}

	view.Render(ctx, writer, "search", page)
}

func (a *searchHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query().Get("query")
	ctx := request.Context()

	articles, err := model.SearchForArticles(query)
	if err != nil {
		renderSearch(ctx, writer, http.StatusInternalServerError, searchData{
			Error: errors.New("Something internal went wrong!").Error(),
		})
		return
	}

	renderSearch(ctx, writer, http.StatusOK, searchData{
		Query:    query,
		Articles: articles,
	})
	return
}

func NewHandler() *searchHandler {
	return &searchHandler{}
}
