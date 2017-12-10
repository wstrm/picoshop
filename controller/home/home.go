package home

import (
	"net/http"
	"log"

	"github.com/willeponken/picoshop/view"
	"github.com/willeponken/picoshop/model"
)

type homeHandler struct {
	http.Handler
}

type homeData struct {
	Error string
	Highlights []model.Article
}

func (h *homeHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.URL.Path != "/" {
		http.NotFound(writer, request)
		return
	}

	highlights, err := model.GetArticleHighlights(10)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	view.Render(request.Context(), writer, "home", view.Page{
		Title: "Picoshop",
		Data: homeData{
			Error: "",
			Highlights: highlights,
		},
	})
}

func NewHandler() *homeHandler {
	return &homeHandler{}
}
