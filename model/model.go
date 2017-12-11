package model

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"github.com/willeponken/picoshop/controller/article"
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

type Order struct {
	Customer   int64
	Address    int64
	Articles   int64
	Status     int
	CreateTime time.Time
}

type Article struct {
	Id          int64
	Name        string
	Description string
	Price       float64
	ImageName   string
	Category    string
	Subcategory string
}

type Category struct {
	Name          string
	Subcategories []Subcategory
}

type Subcategory struct {
	Name     string
	Category string
	Articles []Article
}

type Comment struct {
	Id      int64
	Text    string
	User	int64
}

//go:generate go run $GOPATH/src/github.com/willeponken/picoshop/cmd/inlinesql/main.go -f init.sql -p model -o sql.go
// Open initializes a database connection and forward engineers the ̈́'picoshop' schema with a table setup
func Open(source string) error {
	config, err := mysql.ParseDSN(source)
	if err != nil {
		return err
	}

	config.ParseTime = true

	db, err := sql.Open(driver, config.FormatDSN())
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

func NewArticle(name, description string, price float64, imageName string, category string, subcategory string) Article {
	ensureSubcategoryWithCategory(NewCategory(category), NewSubcategory(subcategory))

	return Article{
		Name:        name,
		Description: description,
		Price:       price,
		ImageName:   imageName,
		Category:    category,
		Subcategory: subcategory,
	}
}

func NewCategory(name string) Category {
	return Category{
		Name: name,
	}
}

func NewSubcategory(name string) Subcategory {
	return Subcategory{
		Name: name,
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

func GetAllOrders() (orders []Order, err error) {
	rows, err := database.Query(`
		SELECT customer, address, articles, status, create_time
		FROM .order`)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		order := Order{}

		err = rows.Scan(
			&order.Customer, &order.Address, &order.Articles,
			&order.Status, &order.CreateTime)
		if err != nil {
			log.Panicln(err)
		}

		orders = append(orders, order)
	}

	err = rows.Err()
	return
}

func SearchForArticles(query string) (articles []Article, err error) {
	rows, err := database.Query(`
		SELECT id, name, description, price, image_name, category, subcategory
		FROM .article WHERE name  LIKE = %_%?%_%`, query)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		article := Article{}

		err = rows.Scan(
			&article.Id, &article.Name, &article.Description,
			&article.Price, &article.ImageName, &article.Category, &article.Subcategory)
		if err != nil {
			log.Panicln(err)
		}

		articles = append(articles, article)
	}

	err = rows.Err()
	return
}

func PutArticle(article Article) (Article, error) {
	result, err := database.Exec(`
		INSERT INTO article
		(name, description, price, image_name, category, subcategory)
		VALUES
		(?, ?, ?, ?, ?)
	`, &article.Name, &article.Description, &article.Price, &article.ImageName, &article.Category, &article.Subcategory)
	if err != nil {
		return Article{}, err
	}

	article.Id, err = result.LastInsertId()
	return article, err
}

func GetAllCategories() (categories []Category, err error) {
	rows, err := database.Query(`
		SELECT (name, subcategories)
		FROM category`)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		category := Category{}

		err = rows.Scan(
			&category.Name, &category.Subcategories)
		if err != nil {
			log.Panicln(err)
		}

		categories = append(categories, category)
	}

	err = rows.Err()
	return
}

func GetSubcategoriesFromCategory(category Category) (subcategories []Subcategory, err error) {
	rows, err := database.Query(`
		SELECT (subcategory.name, subcategory.articles)
		FROM category WHERE category.id=?

		INNER JOIN subcategories
		ON category.subcategories = category_has_subcategories.id

		INNER JOIN subcategory
		ON category_has_subcategories.subcategory = subcategory.id
	`, category.Name)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		subcategory := Subcategory{}

		err = rows.Scan(
			&subcategory.Name, &subcategory.Articles)
		if err != nil {
			log.Panicln(err)
		}

		subcategories = append(subcategories, subcategory)
	}

	err = rows.Err()
	return
}

func GetArticlesFromSubcategory(subcategory Subcategory) (articles []Article, err error) {
	rows, err := database.Query(`
		SELECT (article.id, article.name, article.description, article.price, article.image_name)
		FROM subcategory WHERE subcategory.id=?

		INNER JOIN articles
		ON subcategory.articles = subcategory_has_articles.id

		INNER JOIN article
		ON subcategory_has_articles.article = article.id
	`, subcategory.Name)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		article := Article{}

		err = rows.Scan(
			&article.Id, &article.Name, &article.Description, &article.Price, &article.ImageName)
		if err != nil {
			log.Panicln(err)
		}

		articles = append(articles, article)
	}

	err = rows.Err()
	return
}

func GetArticleHighlights(n uint) (articles []Article, err error) {
	rows, err := database.Query(`
		SELECT id, name, description, price, image_name
		FROM article
		WHERE rand() <= 0.3
		LIMIT ?
	`, n)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		article := Article{}

		err = rows.Scan(
			&article.Id, &article.Name, &article.Description, &article.Price, &article.ImageName)
		if err != nil {
			log.Panicln(err)
		}

		articles = append(articles, article)
	}

	err = rows.Err()
	return
}

func GetArticleById(id int64) (article Article, err error) {
	err = database.QueryRow(`
		SELECT id, name, description, price, image_name
		FROM article
		WHERE id=?
	`, id).Scan(&article.Id, &article.Name, &article.Description, &article.Price, &article.ImageName)

	return
}

func putCategory(category Category) (Category, error) {
	_, err := database.Exec(`
		INSERT IGNORE INTO category
		(name)
		VALUES
		(?)
	`, category.Name)

	return category, err
}

func putSubcategory(subcategory Subcategory) (Subcategory, error) {
	_, err := database.Exec(`
		INSERT IGNORE INTO subcategory
		(name, category)
		VALUES
		(?, ?)
	`, subcategory.Name, subcategory.Category)

	return subcategory, err
}

func addSubcategoryToCategory(category Category, subcategory Subcategory) error {
	_, err := database.Exec(`
		INSERT IGNORE INTO category_has_subcategories
		(category, subcategory)
		VALUES
		(?, ?)
	`, category.Name, subcategory.Name)

	return err
}

func ensureSubcategoryWithCategory(category Category, subcategory Subcategory) {
	cat, err := putCategory(category)
	if err != nil {
		log.Panicln(err)
	}

	subcategory.Category = cat.Name
	subcat, err := putSubcategory(subcategory)
	if err != nil {
		log.Panicln(err)
	}

	err = addSubcategoryToCategory(cat, subcat)
	if err != nil {
		log.Panicln(err)
	}
}

func updateUser(user User)(error){
	_, err := database.Exec(`
		UPDATE user SET Name, Email, PhoneNumber
		(?, ?, ?)`, &user.Name, &user.Email, &user.PhoneNumber)
	return err
}

func addComment(comment Comment)(error) {
	_, err := database.Exec(`
		INSERT INTO comments
		(?, ?)`, &comment.Text, &comment.User)
	return err
}
