package model

import (
	"log"
	"time"
)

type Address struct {
	Id      int64
	Street  string
	CareOf  string
	ZipCode string
	Country string
}

type Order struct {
	Id         int64
	Customer   int64
	Address    Address
	Status     int8
	Articles   []Article
	CreateTime time.Time
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

		// TODO(willeponken): get articles from id in order_has_articles
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

func SetOrderStatus(id int64, status int) error {
	_, err := database.Exec(`
		UPDATE .order SET status = ? WHERE id=?
		`, &status, &id)
	return err
}

/*
func GetOrdersByUserId(id int64) (order []Order, err error) {
	row, err := database.QueryRow(`
		SELECT id, customer, address, status, create_time
		FROM .order
		WHERE
	`)
}
*/
