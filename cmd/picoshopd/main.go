package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/willeponken/picoshop/controller"
)

type flags struct {
	address string
}

var context = flags{
	address: ":8080",
}

func init() {
	flag.StringVar(&context.address, "address", context.address, "Listen address for web server")
}

func main() {
	controller := controller.New()

	log.Printf("Listening on: %s", context.address)
	log.Fatal(http.ListenAndServe(context.address, controller))
}
