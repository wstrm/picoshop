package controller

import (
	"net/http"

	"log"

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

func (a *warehouseHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	orders, err := model.GetAllOrders()
	if err != nil {
		log.Panicln(err)
	}

	view.Render(request.Context(), writer, "warehouse", view.Page{Title: "Warehouse - Picoshop", Data: warehouseData{
		Orders: orders,
	}})
}

func newWarehouseHandler() *warehouseHandler {
	return &warehouseHandler{}
}
