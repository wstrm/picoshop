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
	HasRated bool
}

func (a *articleHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	id := request.URL.Query().Get("id")

	if id != "" {
		id64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

		customer, customerOk := request.Context().Value("Customer").(model.Customer)

		method := request.Method
		switch request.FormValue("_method") {
		case http.MethodPut:
			method = http.MethodPut
		case http.MethodDelete:
			method = http.MethodDelete
		}

		switch method {
		case http.MethodGet: // Just render page
		case http.MethodPost: // Add comment
			if !customerOk {
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			text := request.FormValue("text")
			if text != "" {
				err := model.AddComment(id64, model.Comment{
					Text:     text,
					Customer: customer,
				})
				if err != nil {
					log.Println(err)
					http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
					return
				}
			}
		case http.MethodPut: // Vote
			if !customerOk {
				http.Error(writer, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			vote := request.FormValue("vote")
			var err error
			if vote == "up" {
				err = model.UserRateUp(customer.Id, id64)
			} else if vote == "down" {
				err = model.UserRateDown(customer.Id, id64)
			}

			if err != nil {
				log.Println(err)
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}

		article, err := model.GetArticleById(id64)
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}

		hasRated, _ := model.UserHasRated(customer.Id, id64)

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
				HasRated: hasRated,
			},
		})
		return
	}

	http.NotFound(writer, request)
}

func NewHandler() *articleHandler {
	return &articleHandler{}
}
