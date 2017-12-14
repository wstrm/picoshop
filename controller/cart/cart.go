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

func parseArticleId(article string) (id int64, err error) {
	if err = helper.IsFilled(article); err != nil {
		return
	}

	id, err = strconv.ParseInt(article, 10, 64)
	return
}

func renderCart(writer http.ResponseWriter, request *http.Request, customerId int64) {
	cart, err := model.GetCart(customerId)
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
		renderCart(writer, request, customer.Id)

	case http.MethodPut: // Add article to cart
		articleId, err := parseArticleId(request.FormValue("article"))
		if err != nil {
			log.Println(err)
			http.Error(writer, "Invalid article ID", http.StatusBadRequest)
			return
		}

		err = model.AddArticleToCart(customer.Id, articleId)
		if err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		renderCart(writer, request, customer.Id)

	case http.MethodPost: // Order cart
		address := model.Address{
			Street:  request.FormValue("street"),
			CareOf:  request.FormValue("care_of"),
			ZipCode: request.FormValue("zip_code"),
			Country: request.FormValue("country"),
		}

		err := helper.IsFilled(address.Street, address.CareOf, address.ZipCode, address.Country)
		if err != nil {
			log.Println(err)
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if err := model.OrderCart(customer.Id, address); err != nil {
			log.Println(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Redirect(writer, request, "/order", http.StatusSeeOther)

	case http.MethodDelete:
		article := request.FormValue("article")

		if article != "" { // Delete specific item
			articleId, err := parseArticleId(article)
			if err != nil {
				log.Println(err)
				http.Error(writer, "Invalid article ID", http.StatusBadRequest)
				return
			}

			err = model.DeleteArticleFromCart(customer.Id, articleId)
			if err != nil {
				log.Println(err)
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

		} else {
			err := model.DeleteCart(customer.Id)
			if err != nil {
				log.Println(err)
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}

		renderCart(writer, request, customer.Id)

	default:
		http.Error(writer, "", http.StatusMethodNotAllowed)
	}
}

func NewHandler() *cartHandler {
	return &cartHandler{}
}
