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

type OrderItem struct {
	Article  Article
	Quantity int
}

type Order struct {
	Id         int64
	Customer   int64
	Address    Address
	Status     int64
	Items      []OrderItem
	CreateTime time.Time
}

func newOrderItem(quantity int, article Article) OrderItem {
	return OrderItem{
		Quantity: quantity,
		Article:  article,
	}
}

func GetAllOrders() (orders []Order, err error) {
	tx, err := database.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	orderRows, err := tx.Query(`
		SELECT 	o.id, o.customer, o.status, o.create_time,
			a.street, a.care_of, a.zip_code, a.country
		FROM .order o
		
		INNER JOIN address a
		ON a.id = o.address
		`)
	if err != nil {
		return
	}

	var order Order

	for orderRows.Next() {
		order = Order{}

		err = orderRows.Scan(&order.Id, &order.Customer, &order.Status, &order.CreateTime,
			&order.Address.Street, &order.Address.CareOf, &order.Address.ZipCode,
			&order.Address.Country)
		if err != nil {
			log.Panicln(err)
		}

		orders = append(orders, order)
	}
	err = orderRows.Err()
	if err != nil {
		return
	}
	orderRows.Close()

	articlesStmt, err := tx.Prepare(`
		SELECT a.id, a.name, a.description, a.price, a.image_name, a.category, a.subcategory, o.quantity
		FROM order_has_articles o

		INNER JOIN article a
		ON o.article = a.id

		WHERE o.order = ?
	`)

	for i, _ := range orders {
		articleRows, err := articlesStmt.Query(&orders[i].Id)
		if err != nil {
			log.Panicln(err)
		}

		var quantity int
		var article Article
		var categoryName, subcategoryName string
		for articleRows.Next() {
			err = articleRows.Scan(
				&article.Id, &article.Name, &article.Description, &article.Price,
				&article.ImageName, &categoryName, &subcategoryName, &quantity)
			if err != nil {
				log.Panicln(err)
			}

			article.Category = NewCategory(categoryName)
			article.Subcategory = NewSubcategory(subcategoryName)

			orders[i].Items = append(orders[i].Items, newOrderItem(quantity, article))
		}

		err = articleRows.Err()
		articleRows.Close()

		if err != nil {
			articlesStmt.Close()
			return nil, err
		}
	}

	articlesStmt.Close()

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
