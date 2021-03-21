package model

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-sql-driver/mysql"
	"github.com/wstrm/picoshop/model/forwardengineer"
)

// sql.DB is thread-safe
var database *sql.DB

// driver defines a SQL driver to use
const driver = "mysql"

//go:generate go run $GOPATH/src/github.com/wstrm/picoshop/cmd/inlinesql/main.go -f forwardengineer/schema.sql -p forwardengineer -o forwardengineer/schema.go
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

		os.Exit(0)
	}()
}
