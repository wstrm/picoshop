package controller

import (
	"net/http"
)

type staticHandler struct {
	http.Handler
}

func newStaticHandler() *staticHandler {
	return &staticHandler{
		http.StripPrefix("/static", http.FileServer(http.Dir("static"))),
	}
}
