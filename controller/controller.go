package controller

import (
	"errors"
	"net/http"

	"github.com/willeponken/picoshop/middleware/auth"
)

func internalServerError() error {
	return errors.New("Something internal went wrong!")
}

func invalidFormDataError() error {
	return errors.New("Invalid form data")
}

func New() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", auth.Intercept(newHomeHandler()))
	mux.Handle("/login", auth.Intercept(newLoginHandler()))
	mux.Handle("/register", auth.Intercept(newRegisterHandler()))
	mux.Handle("/user", auth.Intercept(newUserHandler()))
	mux.Handle("/article", auth.Intercept(newArticleHandler()))
	mux.Handle("/cart", auth.Intercept(newCartHandler()))
	mux.Handle("/static/", auth.Intercept(newStaticHandler()))
	mux.Handle("/admin", auth.Protected(newAdminHandler()))
	mux.Handle("/warehouse", auth.Protected(newWarehouseHandler()))
	mux.Handle("/search", auth.Intercept(newSearchHandler()))

	return mux
}
