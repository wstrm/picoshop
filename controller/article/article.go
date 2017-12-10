package article

import (
	"log"
	"net/http"
	"strconv"

	"github.com/willeponken/picoshop/model"
	"github.com/willeponken/picoshop/view"
)

type articleHandler struct {
	http.Handler
}

type articleData struct {
	Error   string
	Article model.Article
}

func (a *articleHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")

	if id != "" {
		id64, err := strconv.ParseInt(id, 10, 64)
		article, err := model.GetArticleById(id64)
		log.Println(err)

		view.Render(request.Context(), writer, "article", view.Page{
			Title: "Article - Picoshop",
			Data: articleData{
				Error:   "",
				Article: article,
			},
		})
		return
	}

	http.NotFound(writer, request)
}

func NewHandler() *articleHandler {
	return &articleHandler{}
}
