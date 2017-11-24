package model

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// sql.DB is thread-safe
var database *sql.DB

// driver defines a SQL driver to use
const driver = "mysql"

const (
	adminRole = iota
	customerRole
	warehouseRole
)

type User struct {
	Id          int64
	Email       string
	Name        string
	Hash        string
	PhoneNumber string
	CreateTime  time.Time
	Addresses   int64
}

type Admin struct {
	User
}

type Warehouse struct {
	User
}

type Customer struct {
	User
	CreditCard int
	Orders     int64
}

//go:generate go run $GOPATH/src/github.com/willeponken/picoshop/cmd/inlinesql/main.go -f init.sql -p model -o sql.go
// Open initializes a database connection and forward engineers the ̈́'picoshop' schema with a table setup
func Open(source string) error {
	db, err := sql.Open(driver, source)
	if err != nil {
		return err
	}
	defer graceful(db.Close)

	database = db

	// Initialize database
	initQueries := getQueries()
	for _, query := range initQueries {
		_, err := database.Exec(query)
		if err != nil {
			return fmt.Errorf("for sql query: %s, got answer: %v", query, err)
		}
	}

	return nil
}

// graceful calls a function upon program exit
func graceful(fn func() error) {
	go func() {
		sig := make(chan os.Signal, 1)
		defer close(sig)

		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig

		err := fn()
		if err != nil {
			panic(err)
		}
	}()
}

func PutUser(user User) (User, error) {
	result, err := database.Exec(`
		INSERT INTO user 
			(email, name, 	hash, 	phone_number, 	create_time)
			VALUES
			(?, 	?,	?, 	?,		?)
		`, user.Email, user.Name, user.Hash, user.PhoneNumber, user.CreateTime)

	if err != nil {
		return User{}, err
	}

	user.Id, _ = result.LastInsertId()

	return user, nil
}

func PutAdmin(admin Admin) (Admin, error) {
	user, err := PutUser(admin.User)
	if err != nil {
		return Admin{}, err
	}

	result, err := database.Exec(`
		INSERT INTO admin
			(user)
			VALUES
			(?)
	`, user.Id)

	log.Println(result)

	return Admin{}, nil //TODO
}

func PutWarehouse(warehouse Warehouse) (Warehouse, error) {
	user, err := PutUser(warehouse.User)
	if err != nil {
		return Warehouse{}, err
	}

	result, err := database.Exec(`
		INSERT INTO warehouse
			(user)
			VALUES
			(?)
	`, user.Id)

	log.Println(result)

	return Warehouse{}, nil //TODO
}

func PutCustomer(customer Customer) (Customer, error) {
	user, err := PutUser(customer.User)
	if err != nil {
		return Customer{}, err
	}

	result, err := database.Exec(`
		INSERT INTO customer
			(user, credit_card)
			VALUES
			(?, ?)
	`, user.Id, customer.CreditCard)

	log.Println(result)

	return Customer{}, nil //TODO
}
