package main

import (
	"flag"
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
