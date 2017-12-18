package model

import "log"

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

	var quantity int
	var article Article

	for rows.Next() {
		err = rows.Scan(
			&article.Id, &article.Name, &article.Description,
			&article.Price, &article.ImageName, &article.Category, &article.Subcategory,
			&quantity)

		if err != nil {
			log.Panicln(err)
		}

		cart.Items = append(cart.Items, newCartItem(quantity, article))
	}

	err = rows.Err()
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

	_, err = tx.Exec(`
		INSERT INTO order_has_articles o
		(o.order, o.article, o.quantity, o.price)
		SELECT ?, c.article, c.quantity, a.price
		FROM cart c
		INNER JOIN article a
		ON c.article = a.id
		WHERE c.customer = ?
	`, &orderId, &customerId)
	if err != nil {
		return
	}

	_, err = tx.Exec(`
		DELETE FROM cart
		WHERE customer = ?
	`, &customerId)
	if err != nil {
		return
	}

	return
}
