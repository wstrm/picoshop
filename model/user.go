package model

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	adminRole = iota
	customerRole
	warehouseRole
)

type User struct {
	Id          int64
	Email       string
	Name        string
	hash        []byte
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

func NewUser(email, name, password, phoneNumber string) User {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), -1)
	if err != nil {
		log.Panicln(err)
	}

	now := time.Now()

	return User{
		Email:       email,
		Name:        name,
		hash:        hash,
		PhoneNumber: phoneNumber,
		CreateTime:  now,
	}
}

func (user User) IsValid() bool {
	if user.Email != "" && user.Name != "" && user.PhoneNumber != "" && !user.CreateTime.IsZero() {
		return true
	}

	return false
}

func NewCustomer(email, name, password, phoneNumber string) Customer {
	return Customer{
		User: NewUser(email, name, password, phoneNumber),
	}
}

func NewAdmin(email, name, password, phoneNumber string) Admin {
	return Admin{
		User: NewUser(email, name, password, phoneNumber),
	}
}

func NewWarehouse(email, name, password, phoneNumber string) Warehouse {
	return Warehouse{
		User: NewUser(email, name, password, phoneNumber),
	}
}

func validPassword(hash []byte, password string) (ok bool) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		return
	}

	ok = true
	return
}

func AuthenticateAdminByEmail(email string, password string) (admin Admin, ok bool) {
	admin, err := GetAdminByEmail(email)
	if err != nil {
		ok = false
		return
	}

	ok = validPassword(admin.hash, password)
	return
}

func AuthenticateCustomerByEmail(email string, password string) (customer Customer, ok bool) {
	customer, err := GetCustomerByEmail(email)
	if err != nil {
		ok = false
		return
	}

	ok = validPassword(customer.hash, password)
	return
}

func AuthenticateWarehouseByEmail(email string, password string) (warehouse Warehouse, ok bool) {
	warehouse, err := GetWarehouseByEmail(email)
	if err != nil {
		ok = false
		return
	}

	ok = validPassword(warehouse.hash, password)
	return
}

func GetAdminByEmail(email string) (admin Admin, err error) {
	err = database.QueryRow(`
		SELECT admin.id, user.id, user.email, user.name, user.hash, user.phone_number, user.create_time
		FROM user
		INNER JOIN admin
		ON user.id = admin.user
		WHERE user.email=LOWER(TRIM(?))
	`, email).Scan(&admin.Id, &admin.User.Id, &admin.User.Email, &admin.User.Name, &admin.User.hash, &admin.User.PhoneNumber, &admin.User.CreateTime)

	return
}

func GetCustomerByEmail(email string) (customer Customer, err error) {
	err = database.QueryRow(`
		SELECT customer.id, user.id, user.email, user.name, user.hash, user.phone_number, user.create_time
		FROM user
		INNER JOIN customer
		ON user.id = customer.user
		WHERE user.email=LOWER(TRIM(?))
	`, email).Scan(&customer.Id, &customer.User.Id, &customer.User.Email, &customer.User.Name, &customer.User.hash, &customer.User.PhoneNumber, &customer.User.CreateTime)

	return
}

func GetWarehouseByEmail(email string) (warehouse Warehouse, err error) {
	err = database.QueryRow(`
		SELECT warehouse.id, user.id, user.email, user.name, user.hash, user.phone_number, user.create_time
		FROM user
		INNER JOIN warehouse
		ON user.id = warehouse.user
		WHERE user.email=LOWER(TRIM(?))
	`, email).Scan(&warehouse.Id, &warehouse.User.Id, &warehouse.User.Email, &warehouse.User.Name, &warehouse.User.hash, &warehouse.User.PhoneNumber, &warehouse.User.CreateTime)

	return
}

func GetUserByEmail(email string) (user User, err error) {
	err = database.QueryRow(`
		SELECT email, name, hash, phone_number, create_time
		FROM user
		WHERE (email=LOWER(TRIM(?)))
	`, email).Scan(&user.Email, &user.Name, &user.hash, &user.PhoneNumber, &user.CreateTime)

	return
}

func PutUser(user User) (User, error) {
	result, err := database.Exec(`
		INSERT INTO user 
			(email, 	name, 	hash, 	phone_number, 	create_time)
			VALUES
			(LOWER(TRIM(?)), 	?,	?, 	?,		?)
		`, user.Email, user.Name, user.hash, user.PhoneNumber, user.CreateTime)

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

func UpdateUser(user User) error {
	_, err := database.Exec(`
		UPDATE .user SET Name, Email, PhoneNumber
		(?, ?, ?)`, &user.Name, &user.Email, &user.PhoneNumber)
	return err
}
