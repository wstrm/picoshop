package controller

import (
	"net/http"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", newHomeHandler())
	mux.Handle("/login", newLoginHandler())
	mux.Handle("/register", newRegisterHandler())
	mux.Handle("/user", newUserHandler())
	mux.Handle("/article", newArticleHandler())
	mux.Handle("/cart", newCartHandler())
	mux.Handle("/static/", newStaticHandler())

	return mux
}
