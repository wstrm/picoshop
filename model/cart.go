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

func DelArticleFromCart(customerId int64, articleId int64) error {
	_, err := database.Exec(`
		START TRANSACTION;

		UPDATE cart
		SET quantity = quantity - 1
		WHERE customer = ? AND article = ?;

		DELETE FROM cart
		WHERE customer = ? AND article = ? AND quantity < 1
		LIMIT 1;

		COMMIT;
	`, &customerId, &articleId, &customerId, &articleId)

	return err
}

func newCartItem(quantity int, article Article) CartItem {
	return CartItem{
		Quantity: quantity,
		Article:  article,
	}
}

func GetCart(customer Customer) (cart Cart, err error) {
	rows, err := database.Query(`
		SELECT (a.id, a.name, a.description, a.price, a.image_name, a.category, a.subcategory,
			c.quantity)
		FROM cart c

		INNER JOIN article a
		ON c.article = a.id

		WHERE c.customer = ?
	`, customer.Id)
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
