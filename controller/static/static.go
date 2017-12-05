package static

import (
	"net/http"
)

type staticHandler struct {
	http.Handler
}

func NewHandler() *staticHandler {
	return &staticHandler{
		http.FileServer(http.Dir("static")),
	}
}
