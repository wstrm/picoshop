package controller

import (
	"net/http"

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

func renderLogin(writer http.ResponseWriter, code int, data interface{}) {
	writer.WriteHeader(code)

	if data == nil {
		data = loginData{}
	}

	page := view.Page{
		Title: "Login - Picoshop",
		Data:  data,
	}

	view.Render(writer, "login", page)
}

func (l *loginHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		renderLogin(writer, http.StatusOK, nil)

	case http.MethodPost:
		err := request.ParseForm()
		if err != nil {
			renderLogin(writer, http.StatusBadRequest, loginData{
				Error: "invalid form data",
			})
			return
		}

		email := request.PostFormValue("email")
		password := request.PostFormValue("password")

		if err := IsFilled(email, password); err != nil {
			renderLogin(writer, http.StatusBadRequest, loginData{
				Error:    err.Error(),
				Email:    email,
				Password: password,
			})
			return
		}

		http.Error(writer, "Not Implemented", http.StatusNotImplemented)

	default:
		http.Error(writer, "Not Allowed", http.StatusMethodNotAllowed)
	}
}

func newLoginHandler() *loginHandler {
	return &loginHandler{}
}
