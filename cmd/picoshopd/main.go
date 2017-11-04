package main

import (
	"net/http"
	"log"

	"github.com/willeponken/picoshop/controller"
)

func main() {
	controller := controller.New()

	log.Printf("Listening on: %s", context.address)
	log.Fatal(http.ListenAndServe(context.address, controller))
}
