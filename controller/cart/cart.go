package cart

import (
	"net/http"

	"github.com/willeponken/picoshop/view"
)

type handler struct {
	http.Handler
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet: // View cart
		view.Render(writer, "cart", view.Page{Title: "Picoshop"})

	case http.MethodPost: // Add article to cart
		http.Error(writer, "", http.StatusNotImplemented)

	case http.MethodPut: // Order items in cart
		http.Error(writer, "", http.StatusNotImplemented)

	case http.MethodDelete:
		pos := request.URL.Query().Get("pos")
		if pos != "" {
			// Delete specific item
		} else {
			// Delete whole cart
		}

	default:
		http.Error(writer, "", http.StatusMethodNotAllowed)
	}
}

func New() *handler {
	return &handler{}
}
