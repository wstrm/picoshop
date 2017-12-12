package cart

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/willeponken/picoshop/controller/helper"
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

func getCustomerFromCtx(ctx context.Context) (customer model.Customer, err error) {
	v := ctx.Value("customer")
	if v == nil {
		err = errors.New("no customer in context")
		return
	}

	customer, ok := v.(model.Customer)
	if !ok {
		err = fmt.Errorf("cannot cast %v to %T", v, model.Customer{})
		return
	}

	if ok := customer.IsValid(); !ok {
		err = fmt.Errorf("customer %v is invalid", customer)
	}

	return
}

func (c *cartHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	customer, err := getCustomerFromCtx(request.Context())
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	switch request.Method {
	case http.MethodGet: // View cart
		view.Render(request.Context(), writer, "cart", view.Page{
			Title: "Picoshop",
			Data:  cartData{},
		})

	case http.MethodPut: // Add article to cart
		id := request.URL.Query().Get("id")
		if err := helper.IsFilled(id); err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		id64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		err = model.AddArticleToCart(customer.Id, id64)
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Error(writer, "Done", http.StatusNotImplemented)

	case http.MethodPost: // Order items in cart
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
