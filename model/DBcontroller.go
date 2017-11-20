package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"log"
)

func main() {
	db, err := sql.Open("mysql",
		"root:toor@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}