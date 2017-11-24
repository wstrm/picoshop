package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/willeponken/picoshop/controller"
	"github.com/willeponken/picoshop/model"
)

type flags struct {
	address string
	source  string
}

var context = flags{
	address: ":8080",
	source:  "root:toor@tcp(127.0.0.1:3306)/picoshopdb",
}

func init() {
	flag.StringVar(&context.address, "address", context.address, "Listen address for web server")
	flag.StringVar(&context.source, "source", context.source, "Database connection source")
	flag.Parse()
}

func main() {
	if err := model.Open(context.source); err != nil {
		log.Fatal(err)
	}

	controller := controller.New()

	log.Printf("Listening on: %s", context.address)
	log.Fatal(http.ListenAndServe(context.address, controller))
}
