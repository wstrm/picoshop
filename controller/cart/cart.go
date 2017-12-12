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

type creditCard struct {
	Number     string
	Expiration string
	Cvc        string
	Holder     string
}

type cartData struct {
	Error      string
	CreditCard creditCard
	Cart       model.Cart
	Address    model.Address
}

func getCustomerFromCtx(ctx context.Context) (customer model.Customer, err error) {
	v := ctx.Value("Customer")
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

	// Because the HTML standard sucks, we fake it till we make it
	method := request.Method
	switch request.FormValue("_method") {
	case http.MethodPut:
		method = http.MethodPut
	case http.MethodDelete:
		method = http.MethodDelete
	}

	switch method {
	case http.MethodGet: // View cart
		view.Render(request.Context(), writer, "cart", view.Page{
			Title: "Picoshop",
			Data:  cartData{},
		})

	case http.MethodPut: // Add article to cart
		id := request.FormValue("article")
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

		cart, err := model.GetCart(customer.Id)
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		view.Render(request.Context(), writer, "cart", view.Page{
			Title: "Picoshop",
			Data: cartData{
				Cart: cart,
			},
		})

	case http.MethodPost: // Order cart
		http.Error(writer, "", http.StatusNotImplemented)

	case http.MethodDelete:
		article := request.FormValue("article")
		if article != "" {
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
