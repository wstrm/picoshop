package warehouse

import (
	"net/http"

	"log"

	"strconv"

	"github.com/wstrm/picoshop/model"
	"github.com/wstrm/picoshop/view"
)

type rootHandler struct {
	http.Handler
}

type orderHandler struct {
	http.Handler
}

type rootData struct {
	Error  string
	Orders []model.Order
}

const (
	pending = iota
	accepted
	shipped
	end
)

func newRootHandler() *rootHandler {
	return &rootHandler{}
}

func newOrderHandler() *orderHandler {
	return &orderHandler{}
}

func (r *rootHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	orders, err := model.GetAllOrders()
	if err != nil {
		log.Panicln(err)
	}

	switch request.Method {
	case http.MethodGet: // View warehouse orders
		view.Render(request.Context(), writer, "warehouse", view.Page{
			Title: "Warehouse - Picoshop",
			Data: rootData{
				Orders: orders,
			},
		})
	default:
		http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (o *orderHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// Because the HTML standard sucks, we fake it till we make it
	method := request.Method
	switch request.FormValue("_method") {
	case http.MethodPut:
		method = http.MethodPut
	case http.MethodDelete:
		method = http.MethodDelete
	}

	id, err := strconv.ParseInt(request.FormValue("id"), 10, 64)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	switch method {
	case http.MethodPost:
		model.SetOrderStatus(id, shipped)

	case http.MethodPut:
		model.SetOrderStatus(id, accepted)

	case http.MethodDelete:
		model.SetOrderStatus(id, end)
	default:
		http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	http.Redirect(writer, request, "/warehouse", http.StatusSeeOther)
}

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", newRootHandler())
	mux.Handle("/order", newOrderHandler())

	return mux
}
