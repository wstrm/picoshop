package model

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-sql-driver/mysql"
	"github.com/willeponken/picoshop/model/forwardengineer"
)

// sql.DB is thread-safe
var database *sql.DB

// driver defines a SQL driver to use
const driver = "mysql"

//go:generate go run $GOPATH/src/github.com/willeponken/picoshop/cmd/inlinesql/main.go -f forwardengineer/schema.sql -p forwardengineer -o forwardengineer/schema.go
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
	initQueries := forwardengineer.GetQueries()
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

func updateUser(user User) error {
	_, err := database.Exec(`
		UPDATE .user SET Name, Email, PhoneNumber
		(?, ?, ?)`, &user.Name, &user.Email, &user.PhoneNumber)
	return err
}

func addComment(comment Comment) error {
	_, err := database.Exec(`
		INSERT INTO comments
		(?, ?)`, &comment.Text, &comment.User)
	return err
}

func SetOrderStatus(id int64, status int) error {
	_, err := database.Exec(`
		UPDATE .order SET status = ? WHERE id=?
		`, &status, &id)
	return err
}

func GetCommentsByArticleId(id int64) (comments []Comment, err error) {
	rows, err := database.Query(`
		SELECT text FROM comments WHERE comments.article = (?)`, &id)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		comment := Comment{}

		err = rows.Scan(
			&comment.Text)
		if err != nil {
			log.Panicln(err)
		}

		comments = append(comments, comment)
	}

	err = rows.Err()
	return
}
