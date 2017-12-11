package cart

import (
	"net/http"

	"github.com/willeponken/picoshop/model"
	"github.com/willeponken/picoshop/view"
)

type cartHandler struct {
	http.Handler
}

type cartData struct {
	Error      string
	CreditCard string
	Expiration string
	Cvc        string
	CardHolder string
	Articles   []model.Article
}

func (c *cartHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet: // View cart
		view.Render(request.Context(), writer, "cart", view.Page{
			Title: "Picoshop",
			Data:  cartData{},
		})

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

func NewHandler() *cartHandler {
	return &cartHandler{}
}
