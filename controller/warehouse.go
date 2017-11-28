package controller

import (
	"net/http"

	"github.com/willeponken/picoshop/model"
	"github.com/willeponken/picoshop/view"
	"log"
)

type warehouseHandler struct {
	http.Handler
}

type warehouseData struct {
	Error  string
	Orders []model.Order
}

func (a *warehouseHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	orders, err := model.GetOrders()
	if err != nil {
		log.Panicln(err)
	}

	view.Render(writer, "warehouse", view.Page{Title: "Warehouse - Picoshop", Data: warehouseData{
		Orders: orders,
	}})
}

func newWarehouseHandler() *warehouseHandler {
	return &warehouseHandler{}
}
