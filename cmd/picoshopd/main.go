package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/wstrm/picoshop/controller"
	"github.com/wstrm/picoshop/model"
)

type flags struct {
	address    string
	source     string
	cert       string
	key        string
	tls        bool
	setupAdmin bool
}

var context = flags{
	address:    ":8080",
	source:     "",
	cert:       "",
	key:        "",
	tls:        false,
	setupAdmin: false,
}

var version = "0.0.1"

func init() {
	flag.StringVar(&context.address, "address", context.address, "Listen address for web server")
	flag.StringVar(&context.source, "source", context.source, "Database connection source")
	flag.StringVar(&context.cert, "cert", context.cert, "Certificate for TLS")
	flag.StringVar(&context.key, "key", context.key, "Key for TLS")
	flag.BoolVar(&context.tls, "tls", context.tls, "Listen using TLS, requires the -cert and -key flags")
	flag.BoolVar(&context.setupAdmin, "setup-admin", context.setupAdmin, "Add a admin user to the database")
	flag.Parse()

	// Log line file:linenumber.
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Prefix log output with "[show (<version>)]".
	log.SetPrefix(fmt.Sprintf("[\033[32m%s\033[0m (%s)] ", "picoshopd", version))

	if context.source == "" {
		log.Fatalln("Please define a MySQL source, example: -source user:password@tcp(127.0.0.1:3306)/picoshop")
	}

	if context.tls && (context.cert == "" || context.key == "") {
		log.Fatalln("Please define both a certificate and key for TLS using the -cert and -key flags")
	}
}

func scanln(s *bufio.Scanner) string {
	s.Scan()
	return s.Text()
}

func setupAdmin() {
	s := bufio.NewScanner(os.Stdin)

	fmt.Print("Name: ")
	name := scanln(s)

	fmt.Print("E-mail: ")
	email := scanln(s)

	fmt.Print("Password: \033[8m") // hide input
	password := scanln(s)
	fmt.Print("\033[28m") // show input

	fmt.Print("Phone number: ")
	phoneNumber := scanln(s)

	admin := model.NewAdmin(email, name, password, phoneNumber)
	admin, err := model.PutAdmin(admin)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf(`
New admin user added:
	ID: %d
	Name: %s
	E-mail: %s
	Phone number: %s
	`, admin.Id, name, email, phoneNumber)
}

func main() {
	if err := model.Open(context.source); err != nil {
		log.Fatal(err)
	}

	if context.setupAdmin {
		setupAdmin()
		os.Exit(0)
	}

	controller := controller.New()

	log.Printf("Listening on: %s", context.address)
	if context.tls {
		log.Fatal(http.ListenAndServeTLS(context.address, context.cert, context.key, controller))
	} else {
		log.Fatal(http.ListenAndServe(context.address, controller))
	}
}
