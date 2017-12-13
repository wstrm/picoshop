package warehouse

import (
	"net/http"

	"log"

	"strconv"

	"github.com/willeponken/picoshop/model"
	"github.com/willeponken/picoshop/view"
)

type warehouseHandler struct {
	http.Handler
}

type warehouseData struct {
	Error  string
	Orders []model.Order
}

const (
	pending = iota
	accepted
	shipped
	end
)

func (a *warehouseHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	orders, err := model.GetAllOrders()
	if err != nil {
		log.Panicln(err)
	}

	switch request.Method {
	case http.MethodGet: //view warehouse orders

	case http.MethodPost:
		id, _ := strconv.ParseInt(request.FormValue("id"), 10, 64)
		model.SetOrderStatus(id, shipped)

	case http.MethodPut:
		id, _ := strconv.ParseInt(request.FormValue("id"), 10, 64)
		model.SetOrderStatus(id, accepted)

	case http.MethodDelete:
		id, _ := strconv.ParseInt(request.FormValue("id"), 10, 64)
		model.SetOrderStatus(id, end)
	}

	view.Render(request.Context(), writer, "warehouse", view.Page{Title: "Warehouse - Picoshop", Data: warehouseData{
		Orders: orders,
	}})
}

func NewHandler() *warehouseHandler {
	return &warehouseHandler{}
}
