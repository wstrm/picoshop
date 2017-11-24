package model

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
)

// sql.DB is thread-safe
var database *sql.DB

// driver defines a SQL driver to use
const driver = "mysql"

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
