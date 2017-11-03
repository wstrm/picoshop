package controller

import (
	"net/http"

	"github.com/willeponken/picoshop/controller/home"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", home.New())

	return mux
}
