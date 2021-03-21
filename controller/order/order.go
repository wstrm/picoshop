package order

import (
	"log"
	"net/http"

	"github.com/wstrm/picoshop/model"
	"github.com/wstrm/picoshop/view"
)

type orderHandler struct {
	http.Handler
}

type orderData struct {
	Error  string
	Orders []model.Order
}

func (o *orderHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		ctx := request.Context()

		c := ctx.Value("Customer")
		if c == "" {
			log.Println("no customer")
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		customer, ok := c.(model.Customer)
		if !ok {
			log.Println("cannot cast")
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// TODO(wstrm): query only user specific orders
		allOrders, err := model.GetAllOrders()
		if err != nil {
			log.Panicln(err)
		}

		var data orderData
		for _, order := range allOrders {
			if order.Customer == customer.Id {
				data.Orders = append(data.Orders, order)
			}
		}

		view.Render(request.Context(), writer, "order", view.Page{
			Title: "Orders - Picoshop",
			Data:  data,
		})

	default:
		http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func NewHandler() *orderHandler {
	return &orderHandler{}
}
