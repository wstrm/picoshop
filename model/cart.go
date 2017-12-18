package model

import (
	"database/sql"
	"log"
)

type CartItem struct {
	Quantity int
	Article  Article
}

type Cart struct {
	Items []CartItem
}

func AddArticleToCart(customerId int64, articleId int64) error {
	_, err := database.Exec(`
		INSERT INTO cart
		(customer, article, quantity)
		VALUES
		(?, ?, 1)

		ON DUPLICATE KEY UPDATE
		quantity = quantity + 1
	`, &customerId, &articleId)

	return err
}

func DeleteArticleFromCart(customerId int64, articleId int64) (err error) {
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

	update, err := tx.Prepare(`
		UPDATE cart
		SET quantity = (quantity - 1)
		WHERE (customer = ?) AND (article = ?);
	`)
	if err != nil {
		return
	}

	delete, err := tx.Prepare(`
		DELETE FROM cart
		WHERE (customer = ?) AND (article = ?) AND (quantity < 1)
		LIMIT 1;
	`)
	if err != nil {
		return
	}

	if _, err = update.Exec(&customerId, &articleId); err != nil {
		return
	}
	if _, err = delete.Exec(&customerId, &articleId); err != nil {
		return
	}

	return
}

func DeleteCart(customerId int64) error {
	_, err := database.Exec(`
		DELETE FROM cart
		WHERE customer = ?
	`, &customerId)

	return err
}

func newCartItem(quantity int, article Article) CartItem {
	return CartItem{
		Quantity: quantity,
		Article:  article,
	}
}

func scanCart(rows *sql.Rows) (cart Cart, err error) {
	var quantity int
	var article Article
	var categoryName, subcategoryName string

	for rows.Next() {
		err = rows.Scan(
			&article.Id, &article.Name, &article.Description,
			&article.Price, &article.ImageName, &categoryName, &subcategoryName,
			&quantity)

		if err != nil {
			log.Panicln(err)
		}

		article.Category = NewCategory(categoryName)
		article.Subcategory = NewSubcategory(subcategoryName)

		cart.Items = append(cart.Items, newCartItem(quantity, article))
	}

	err = rows.Err()
	return
}

func GetCart(customerId int64) (cart Cart, err error) {
	rows, err := database.Query(`
		SELECT a.id, a.name, a.description, a.price, a.image_name, a.category, a.subcategory, c.quantity
		FROM cart c

		INNER JOIN article a
		ON c.article = a.id

		WHERE c.customer = ?
	`, customerId)
	if err != nil {
		return
	}

	defer rows.Close()

	cart, err = scanCart(rows)

	return
}

func OrderCart(customerId int64, address Address) (err error) {
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

	// Create a new address
	addressResult, err := tx.Exec(`
		INSERT INTO address
		(street, care_of, zip_code, country, customer)
		VALUES
		(?, ?, ?, ?, ?)
	`, &address.Street, &address.CareOf, &address.ZipCode, &address.Country, &customerId)
	if err != nil {
		return
	}

	addressId, err := addressResult.LastInsertId()
	if err != nil {
		return
	}

	// Create new order for customer with address
	orderResult, err := tx.Exec(`
		INSERT INTO .order
		(customer, address)
		VALUES
		(?, ?)
	`, &customerId, &addressId)
	if err != nil {
		return
	}

	orderId, err := orderResult.LastInsertId()
	if err != nil {
		return
	}

	// Insert all articles in cart to order_has_articles
	_, err = tx.Exec(`
		INSERT INTO order_has_articles
		(order_has_articles.order, article, quantity, price)
		SELECT ?, c.article, c.quantity, a.price
		FROM cart c
		INNER JOIN article a
		ON c.article = a.id
		WHERE c.customer = ?
		AND a.in_stock > 0
	`, &orderId, &customerId)
	if err != nil {
		return
	}

	// Decrement in stock for each article in order
	_, err = tx.Exec(`
		UPDATE article a
		LEFT JOIN order_has_articles o
		ON a.id = o.article AND o.order = ?
		SET a.in_stock = a.in_stock - o.quantity
	`, &orderId)
	if err != nil {
		return
	}

	// Delete cart (so user can start all over again)
	_, err = tx.Exec(`
		DELETE FROM cart
		WHERE customer = ?
	`, &customerId)
	if err != nil {
		return
	}

	return
}
