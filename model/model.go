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
	"golang.org/x/crypto/bcrypt"
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
	Hash        []byte
	PhoneNumber string
	CreateTime  time.Time
	Addresses   int64
}

type Admin struct {
	User
	Id int64
}

type Warehouse struct {
	User
	Id int64
}

type Customer struct {
	User
	Id     int64
	Orders int64
}

type Order struct {
	Customer   int
	Address    int
	Articles   int
	Status     int
	CreateTime time.Time
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

func NewUser(email, name, password, phoneNumber string) User {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), -1)
	if err != nil {
		log.Panicln(err)
	}

	now := time.Now()

	return User{
		Email:       email,
		Name:        name,
		Hash:        hash,
		PhoneNumber: phoneNumber,
		CreateTime:  now,
	}
}

func NewCustomer(email, name, password, phoneNumber string) Customer {
	return Customer{
		User: NewUser(email, name, password, phoneNumber),
	}
}

func ValidPassword(email, password string) bool {
	hash, err := getUserHash(email)
	if err != nil {
		return false
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return false
	}

	return true
}

func getUserHash(email string) ([]byte, error) {
	var hash []byte

	err := database.QueryRow(`
		SELECT (hash) FROM user WHERE (email=LOWER(TRIM(?)))
	`, email).Scan(&hash)

	if err != nil {
		return nil, err
	}

	return hash, nil
}

func PutUser(user User) (User, error) {
	result, err := database.Exec(`
		INSERT INTO user 
			(email, 	name, 	hash, 	phone_number, 	create_time)
			VALUES
			(LOWER(TRIM(?)), 	?,	?, 	?,		?)
		`, user.Email, user.Name, user.Hash, user.PhoneNumber, user.CreateTime)

	if err != nil {
		return User{}, err
	}

	user.Id, err = result.LastInsertId()
	if err != nil {
		return User{}, err
	}

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
	if err != nil {
		return Admin{}, err
	}

	admin.Id, err = result.LastInsertId()
	if err != nil {
		return Admin{}, err
	}

	return admin, nil
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
	if err != nil {
		return Warehouse{}, err
	}

	warehouse.Id, err = result.LastInsertId()
	if err != nil {
		return Warehouse{}, err
	}

	return warehouse, nil
}

func PutCustomer(customer Customer) (Customer, error) {
	user, err := PutUser(customer.User)
	if err != nil {
		return Customer{}, err
	}

	result, err := database.Exec(`
		INSERT INTO customer
			(user)
			VALUES
			(?)
	`, user.Id)
	if err != nil {
		return Customer{}, err
	}

	customer.Id, err = result.LastInsertId()
	if err != nil {
		return Customer{}, err
	}

	return customer, nil
}

func GetOrders() (orders []Order, err error) {
	rows, err := database.Query(`
		SELECT customer, address, articles, status, create_time FROM .order`)
	defer rows.Close()

	if err != nil {
		return
	}

	for rows.Next() {
		order := Order{}

		err = rows.Scan(&order.Customer, &order.Address, &order.Articles)
		if err != nil {
			log.Panicln(err)
		}

		orders = append(orders, order)
	}

	err = rows.Err()
	if err != nil {
		return
	}

	return
}
