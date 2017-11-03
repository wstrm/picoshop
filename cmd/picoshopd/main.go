package main

import (
	"net/http"
	"log"

	"github.com/willeponken/picoshop/controller"
)

func main() {
	controller := controller.New()

	log.Fatal(http.ListenAndServe(":8080", controller))
}
