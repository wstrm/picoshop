package category

import (
	"context"
	"log"
	"net/http"

	"github.com/wstrm/picoshop/model"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		categories, err := model.GetAllCategories()

		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(request.Context(), "Categories", categories)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
