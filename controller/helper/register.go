package helper

import (
	"errors"
	"net/http"

	"github.com/willeponken/picoshop/middleware/auth"
	"github.com/willeponken/picoshop/model"
)

type RegisterForm struct {
	Email          string
	Name           string
	PhoneNumber    string
	Password       string
	PasswordRetype string
	Type           string
}

type RegisterResult struct {
	Error          string
	Email          string
	Name           string
	PhoneNumber    string
	Password       string
	PasswordRetype string
	Type           string
}

func LegalPassword(password, passwordRetype string) error {
	if password != passwordRetype {
		return errors.New("re-typed password must equal original")
	}

	if len(password) < 8 {
		return errors.New("password must be atleast 8 characters long")
	}

	return nil
}

func ParseRegisterFormValues(request *http.Request) RegisterForm {
	return RegisterForm{
		Email:          request.PostFormValue("email"),
		Name:           request.PostFormValue("name"),
		PhoneNumber:    request.PostFormValue("phone-number"),
		Password:       request.PostFormValue("password"),
		PasswordRetype: request.PostFormValue("password-retype"),
	}
}

func Register(userType string, form RegisterForm) (user auth.User, result RegisterResult, code int) {
	if err := IsFilled(form.Email, form.Name, form.PhoneNumber, form.Password, form.PasswordRetype); err != nil {
		result = RegisterResult{
			Error:          err.Error(),
			Email:          form.Email,
			Name:           form.Name,
			PhoneNumber:    form.PhoneNumber,
			Password:       form.Password,
			PasswordRetype: form.PasswordRetype,
			Type:           userType,
		}
		code = http.StatusBadRequest
		return
	}

	if err := LegalPassword(form.Password, form.PasswordRetype); err != nil {
		result = RegisterResult{
			Error:       err.Error(),
			Email:       form.Email,
			Name:        form.Name,
			PhoneNumber: form.PhoneNumber,
			Type:        userType,
			// Ignore passwords (force the user to retype invalid data)
		}
		code = http.StatusBadRequest
		return
	}

	var err error
	switch userType {
	case "admin":
		user, err = model.PutAdmin(model.NewAdmin(form.Email, form.Name, form.Password, form.PhoneNumber))
	case "warehouse":
		user, err = model.PutWarehouse(model.NewWarehouse(form.Email, form.Name, form.Password, form.PhoneNumber))
	default: // "customer"
		user, err = model.PutCustomer(model.NewCustomer(form.Email, form.Name, form.Password, form.PhoneNumber))
	}

	if err != nil {
		sqlCode := GetSqlErrorCode(err)
		var userErr error

		switch sqlCode {
		case DuplicateKeySqlError:
			userErr = EmailAlreadyRegisteredError(form.Email)
		default:
			userErr = InternalServerError()
		}

		result = RegisterResult{
			Error: userErr.Error(),
			Type:  userType,
			// Ignore data, apparently it's dangerous!
		}
		code = http.StatusInternalServerError
		return
	}

	code = http.StatusOK
	return
}
