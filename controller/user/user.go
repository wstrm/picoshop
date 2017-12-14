package user

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/willeponken/picoshop/controller/helper"
	"github.com/willeponken/picoshop/middleware/auth"
	"github.com/willeponken/picoshop/model"
	"github.com/willeponken/picoshop/view"
)

type userHandler struct {
	http.Handler
	auth *auth.Manager
}

type userData struct {
	User model.User
}

func renderUser(ctx context.Context, writer http.ResponseWriter, code int, data interface{}) {
	writer.WriteHeader(code)

	if data == nil {
		data = userData{}
	}

	page := view.Page{
		Title: "User - Picoshop",
		Data:  data,
	}

	view.Render(ctx, writer, "user", page)
}

func (u *userHandler) getUser(request *http.Request) (user model.User, err error) {
	authUser, err := u.auth.GetUser(request)
	if err != nil {
		return
	}

	var ok bool
	switch authUser.(type) {
	case model.Admin:
		var a model.Admin
		a, ok = authUser.(model.Admin)
		if ok {
			user = a.User
		}
	case model.Warehouse:
		var w model.Warehouse
		w, ok = authUser.(model.Warehouse)
		if ok {
			user = w.User
		}
	case model.Customer:
		var c model.Customer
		c, ok = authUser.(model.Customer)
		if ok {
			user = c.User
		}
	default:
		err = fmt.Errorf("unsupported user type: %T", authUser)
		return
	}

	if !ok {
		err = fmt.Errorf("cannot type assert %v to model.User", authUser)
	}

	return
}

func (u *userHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	user, err := u.getUser(request)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	switch request.Method {
	case http.MethodGet:
		renderUser(request.Context(), writer, http.StatusOK, userData{user})
		return

	case http.MethodPost:
		name := request.FormValue("name")
		phone := request.FormValue("phone")
		if err := helper.IsFilled(name, phone); err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		user.Name = name
		user.PhoneNumber = phone

		if err := model.UpdateUser(user); err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		u.auth.Logout(writer, request)
		http.Redirect(writer, request, "/", http.StatusSeeOther)

	default:
		http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func NewHandler(auth *auth.Manager) *userHandler {
	return &userHandler{
		auth: auth,
	}
}
