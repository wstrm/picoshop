package controller

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/willeponken/picoshop/model"
	"github.com/willeponken/picoshop/view"
)

type registerHandler struct {
	http.Handler
}

type registerData struct {
	Error          string
	Email          string
	Name           string
	PhoneNumber    string
	Password       string
	PasswordRetype string
}

func renderRegister(writer http.ResponseWriter, code int, data interface{}) {
	writer.WriteHeader(code)

	if data == nil {
		data = registerData{}
	}

	page := view.Page{
		Title: "Register - Picoshop",
		Data:  data,
	}

	view.Render(writer, "register", page)
}

func legalPassword(password, passwordRetype string) error {
	if password != passwordRetype {
		return errors.New("re-typed password must equal original")
	}

	if len(password) < 8 {
		return errors.New("password must be atleast 8 characters long")
	}

	return nil
}

func (r *registerHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {

	case http.MethodGet: // Serve register view
		renderRegister(writer, http.StatusOK, nil)

	case http.MethodPost: // Retreive user registration
		err := request.ParseForm()
		if err != nil {
			renderRegister(writer, http.StatusBadRequest, registerData{
				Error: "invalid form data",
			})
			return
		}

		email := request.PostFormValue("email")
		name := request.PostFormValue("name")
		phoneNumber := request.PostFormValue("phone-number")
		password := request.PostFormValue("password")
		passwordRetype := request.PostFormValue("password-retype")

		if err := IsFilled(email, name, phoneNumber, password, passwordRetype); err != nil {
			renderRegister(writer, http.StatusBadRequest, registerData{
				Error:          err.Error(),
				Email:          email,
				Name:           name,
				PhoneNumber:    phoneNumber,
				Password:       password,
				PasswordRetype: passwordRetype,
			})
			return
		}

		if err := legalPassword(password, passwordRetype); err != nil {
			renderRegister(writer, http.StatusBadRequest, registerData{
				Error:       err.Error(),
				Email:       email,
				Name:        name,
				PhoneNumber: phoneNumber,
				// Ignore passwords (force the user to retype invalid data)
			})
			return
		}

		_, err = model.PutCustomer(model.NewCustomer(email, name, password, phoneNumber))
		if err != nil {
			code := getSqlErrorCode(err)
			var userErr error

			switch code {
			case DuplicateKeySqlError:
				userErr = fmt.Errorf("The email address '%s' is already registered", email)
			default:
				userErr = errors.New("Something internal went wrong!")
			}

			renderRegister(writer, http.StatusInternalServerError, registerData{
				Error: userErr.Error(),
				// Ignore data, apparently it's dangerous!
			})
			return
		}

		http.Redirect(writer, request, "/", http.StatusSeeOther) // See RFC 2616 (redirect after POST)
	}
}

func newRegisterHandler() *registerHandler {
	return &registerHandler{}
}
