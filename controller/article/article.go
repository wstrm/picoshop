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
	Error    string
	Article  model.Article
	Comments []model.Comment
}

func (a *articleHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")

	if id != "" {
		id64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		article, err := model.GetArticleById(id64)
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}

		switch request.Method {
		case http.MethodGet: // Just render page
		case http.MethodPost: // Add comment
			customer, ok := request.Context().Value("Customer").(model.Customer)
			if !ok {
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			text := request.FormValue("text")
			if text != "" {
				err := model.AddComment(id64, model.Comment{
					Text:     text,
					Customer: customer.Id,
				})
				if err != nil {
					log.Println(err)
					http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
			}
		}

		comments, err := model.GetCommentsByArticleId(id64)
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		view.Render(request.Context(), writer, "article", view.Page{
			Title: "Article - Picoshop",
			Data: articleData{
				Error:    "",
				Article:  article,
				Comments: comments,
			},
		})
		return
	}

	http.NotFound(writer, request)
}

func NewHandler() *articleHandler {
	return &articleHandler{}
}
