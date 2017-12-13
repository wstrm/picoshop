package warehouse

import (
	"net/http"

	"log"

	"github.com/willeponken/picoshop/model"
	"github.com/willeponken/picoshop/view"
	"strconv"
)

type warehouseHandler struct {
	http.Handler
}

type warehouseData struct {
	Error  string
	Orders []model.Order
}

const pending = 0
const accepted = 1
const shipped = 2
const end = 3
func (a *warehouseHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	orders, err := model.GetAllOrders()
	if err != nil {
		log.Panicln(err)
	}
	switch request.Method {
	case http.MethodGet: //view warehouse orders
		view.Render(request.Context(), writer, "warehouse", view.Page{Title: "Warehouse - Picoshop", Data: warehouseData{
		Orders: orders,
	}})
	case http.MethodPost:
		id, _ := strconv.ParseInt(request.FormValue("id"), 10, 64)
		model.SetOrderStatus(id, shipped)
		view.Render(request.Context(), writer, "warehouse", view.Page{Title: "Warehouse - Picoshop", Data: warehouseData{
			Orders: orders,
		}})
	case http.MethodPut:
		id, _ := strconv.ParseInt(request.FormValue("id"), 10, 64)
		model.SetOrderStatus(id, accepted)
		view.Render(request.Context(), writer, "warehouse", view.Page{Title: "Warehouse - Picoshop", Data: warehouseData{
			Orders: orders,
		}})

	case http.MethodDelete:
		id, _ := strconv.ParseInt(request.FormValue("id"), 10, 64)
		model.SetOrderStatus(id, end)
		view.Render(request.Context(), writer, "warehouse", view.Page{Title: "Warehouse - Picoshop", Data: warehouseData{
			Orders: orders,
		}})
		}
}

func NewHandler() *warehouseHandler {
	return &warehouseHandler{}
}

func (c *cartHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet: // View cart
		view.Render(request.Context(), writer, "cart", view.Page{Title: "Picoshop"})

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