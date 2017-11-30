package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/willeponken/picoshop/middleware/auth"
	"github.com/willeponken/picoshop/model"
	"github.com/willeponken/picoshop/view"
)

type loginHandler struct {
	http.Handler
}

type loginData struct {
	Error    string
	Email    string
	Password string
}

func invalidLoginCredentialsError() error {
	return errors.New("Invalid login credentials")
}

func renderLogin(ctx context.Context, writer http.ResponseWriter, code int, data interface{}) {
	writer.WriteHeader(code)

	if data == nil {
		data = loginData{}
	}

	page := view.Page{
		Title: "Login - Picoshop",
		Data:  data,
	}

	view.Render(ctx, writer, "login", page)
}

func (l *loginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	switch request.Method {
	case http.MethodGet:
		renderLogin(ctx, writer, http.StatusOK, nil)

	case http.MethodPost:
		err := request.ParseForm()
		if err != nil {
			renderLogin(ctx, writer, http.StatusBadRequest, loginData{
				Error: invalidFormDataError().Error(),
			})
			return
		}

		email := request.PostFormValue("email")
		password := request.PostFormValue("password")

		if err := IsFilled(email, password); err != nil {
			renderLogin(ctx, writer, http.StatusBadRequest, loginData{
				Error:    err.Error(),
				Email:    email,
				Password: password,
			})
			return
		}

		if model.ValidPassword(email, password) {
			err := auth.Login(email, writer, request)
			if err != nil {
				renderLogin(ctx, writer, http.StatusInternalServerError, registerData{
					Error: internalServerError().Error(),
				})
			}

			http.Redirect(writer, request, "/", http.StatusSeeOther)
		} else {
			renderLogin(ctx, writer, http.StatusUnauthorized, loginData{
				Error: invalidLoginCredentialsError().Error(),
			})
		}

	default:
		http.Error(writer, "Not Allowed", http.StatusMethodNotAllowed)
	}
}

func newLoginHandler() *loginHandler {
	return &loginHandler{}
}
