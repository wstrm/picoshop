package category

import (
	"log"
	"net/http"
	"strings"

	"github.com/wstrm/picoshop/model"
	"github.com/wstrm/picoshop/view"
)

type categoryHandler struct {
	http.Handler
}

type categoryData struct {
	Error       string
	Category    string
	Subcategory string
	Articles    []model.Article
}

func (c *categoryHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	path := strings.SplitN(strings.Trim(request.URL.Path, "/"), "/", 2)

	var (
		articles              []model.Article
		err                   error
		category, subcategory string
	)

	switch len(path) {
	case 1: // Get all articles for category
		category = path[0]
		articles, err = model.GetArticlesByCategory(category)
	case 2: // Get all articles for subcategory
		category = path[0]
		subcategory = path[1]
		articles, err = model.GetArticlesBySubcategory(subcategory)
	default:
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	view.Render(request.Context(), writer, "category", view.Page{
		Title: "Category - Picoshop",
		Data: categoryData{
			Error:       "",
			Category:    category,
			Subcategory: subcategory,
			Articles:    articles,
		},
	})
}

func NewHandler() *categoryHandler {
	return &categoryHandler{}
}
