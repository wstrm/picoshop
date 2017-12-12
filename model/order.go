package model

import (
	"log"
	"time"
)

type Order struct {
	Id         int64
	Customer   int64
	Address    int64
	Status     int8
	CreateTime time.Time
}

type Address struct {
	Id      int64
	Street  string
	CareOf  string
	ZipCode string
	Country string
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
