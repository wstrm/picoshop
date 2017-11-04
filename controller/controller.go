package controller

import (
	"net/http"

	"github.com/willeponken/picoshop/controller/home"
	"github.com/willeponken/picoshop/controller/login"
	"github.com/willeponken/picoshop/controller/register"
	"github.com/willeponken/picoshop/controller/user"
	"github.com/willeponken/picoshop/controller/article"
	"github.com/willeponken/picoshop/controller/cart"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", home.New())
	mux.Handle("/login", login.New())
	mux.Handle("/register", register.New())
	mux.Handle("/user", user.New())
	mux.Handle("/article", article.New())
	mux.Handle("/cart", cart.New())

	return mux
}
