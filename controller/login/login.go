package login

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/wstrm/picoshop/controller/helper"
	"github.com/wstrm/picoshop/middleware/auth"
	"github.com/wstrm/picoshop/model"
	"github.com/wstrm/picoshop/view"
)

type loginHandler struct {
	http.Handler
	authManager *auth.Manager
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
				Error: helper.InvalidFormDataError().Error(),
			})
			return
		}

		email := request.PostFormValue("email")
		password := request.PostFormValue("password")
		userType := request.PostFormValue("type")

		if err := helper.IsFilled(email, password); err != nil {
			renderLogin(ctx, writer, http.StatusBadRequest, loginData{
				Error:    err.Error(),
				Email:    email,
				Password: password,
			})
			return
		}

		var (
			user auth.User
			ok   bool
		)
		switch userType {
		case "admin":
			user, ok = model.AuthenticateAdminByEmail(email, password)
		case "customer":
			user, ok = model.AuthenticateCustomerByEmail(email, password)
		case "warehouse":
			user, ok = model.AuthenticateWarehouseByEmail(email, password)
		default:
			renderLogin(ctx, writer, http.StatusBadRequest, loginData{
				Error: errors.New("Missing user type field").Error(),
			})
			return
		}

		if ok {
			err := l.authManager.Login(user, writer, request)
			if err != nil {
				log.Println(err)

				renderLogin(ctx, writer, http.StatusInternalServerError, loginData{
					Error: helper.InternalServerError().Error(),
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

func NewHandler(authManager *auth.Manager) *loginHandler {
	return &loginHandler{
		authManager: authManager,
	}
}
