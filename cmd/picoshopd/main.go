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
	cert    string
	key     string
	tls     bool
}

var context = flags{
	address: ":8080",
	source:  "",
	cert:    "",
	key:     "",
	tls:     false,
}

func init() {
	flag.StringVar(&context.address, "address", context.address, "Listen address for web server")
	flag.StringVar(&context.source, "source", context.source, "Database connection source")
	flag.StringVar(&context.cert, "cert", context.cert, "Certificate for TLS")
	flag.StringVar(&context.key, "key", context.key, "Key for TLS")
	flag.BoolVar(&context.tls, "tls", context.tls, "Listen using TLS, requires the -cert and -key flags")
	flag.Parse()

	if context.source == "" {
		log.Fatalln("Please define a MySQL source, example: -source user:password@tcp(127.0.0.1:3306)/picoshop")
	}

	if context.tls && (context.cert == "" || context.key == "") {
		log.Fatalln("Please define both a certificate and key for TLS using the -cert and -key flags")
	}
}

func main() {
	if err := model.Open(context.source); err != nil {
		log.Fatal(err)
	}

	controller := controller.New()

	log.Printf("Listening on: %s", context.address)
	if context.tls {
		log.Fatal(http.ListenAndServeTLS(context.address, context.cert, context.key, controller))
	} else {
		log.Fatal(http.ListenAndServe(context.address, controller))
	}
}
