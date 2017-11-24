package controller

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/willeponken/picoshop/model"
	"github.com/willeponken/picoshop/view"
)

type registerHandler struct {
	http.Handler
}

type registerData struct {
	Error          string
	Email          string
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

func validatePassword(password, passwordRetype string) error {
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
		password := request.PostFormValue("password")
		passwordRetype := request.PostFormValue("password-retype")

		if err := IsFilled(email, password, passwordRetype); err != nil {
			renderRegister(writer, http.StatusBadRequest, registerData{
				Error:          err.Error(),
				Email:          email,
				Password:       password,
				PasswordRetype: passwordRetype,
			})
			return
		}

		if err := validatePassword(password, passwordRetype); err != nil {
			renderRegister(writer, http.StatusBadRequest, registerData{
				Error: err.Error(),
				Email: email,
				// Ignore passwords (force the user to retype invalid data)
			})
			return
		}

		customer, err := model.PutCustomer(model.Customer{
			User: model.User{
				Email:       email,
				Name:        "A name",         //TODO implement name on client side
				Hash:        password,         //TODO bcrypt
				PhoneNumber: "A phone number", //TODO see above
				CreateTime:  time.Now(),
			},
			CreditCard: 0, //TODO credit inster
		})

		if err != nil {
			log.Println(err)
			renderRegister(writer, http.StatusInternalServerError, registerData{
				Error: errors.New("Something internal went wrong!").Error(),
				// Ignore data, apparently it's dangerous!
			})
			return
		}

		log.Println(customer)

		http.Redirect(writer, request, "/", http.StatusSeeOther) // See RFC 2616 (redirect after POST)
	}
}

func newRegisterHandler() *registerHandler {
	return &registerHandler{}
}
