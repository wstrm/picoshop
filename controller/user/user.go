package user

import (
	"net/http"

	"github.com/willeponken/picoshop/model"
	"github.com/willeponken/picoshop/view"
)

type userHandler struct {
	http.Handler
}

type userData struct {
	User model.User
}

func (u *userHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	admin := ctx.Value("Admin")
	warehouse := ctx.Value("Warehouse")
	customer := ctx.Value("Customer")

	var user model.User
	if admin != "" {
		a, ok := admin.(model.Admin)
		user = a.User
		if !ok {
			http.Error(writer, "I fucked up", http.StatusInternalServerError)
			return
		}
	} else if warehouse != "" {
		w, ok := warehouse.(model.Warehouse)
		user = w.User
		if !ok {
			http.Error(writer, "I fucked up", http.StatusInternalServerError)
			return
		}
	} else if customer != "" {
		c, ok := customer.(model.Customer)
		user = c.User
		if !ok {
			http.Error(writer, "I fucked up", http.StatusInternalServerError)
			return
		}
	}

	switch request.Method {
	case http.MethodGet:
		view.Render(request.Context(), writer, "user", view.Page{
			Title: "Picoshop",
			Data: userData{
				User: user,
			},
		})

	case http.MethodPost:
		http.Error(writer, "", http.StatusNotImplemented)

	default:
		http.Error(writer, "", http.StatusMethodNotAllowed)
	}
}

func NewHandler() *userHandler {
	return &userHandler{}
}
