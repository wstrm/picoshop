package model

import (
	"log"
)

type Comment struct {
	Id       int64
	Text     string
	Customer Customer
}

func GetCommentsByArticleId(id int64) (comments []Comment, err error) {
	rows, err := database.Query(`
		SELECT comment.text, comment.id, customer.id, user.id, user.email, user.name,
		       user.phone_number, user.create_time
		FROM comment

		INNER JOIN customer
		ON customer.id = comment.customer

		INNER JOIN user
		ON user.id = customer.user

		WHERE comment.article = ?
	`, &id)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		comment := Comment{}
		customer := Customer{}

		err = rows.Scan(
			&comment.Text, &comment.Id, &customer.Id, &customer.User.Id, &customer.User.Email,
			&customer.User.Name, &customer.User.PhoneNumber, &customer.User.CreateTime)
		if err != nil {
			log.Panicln(err)
		}

		comment.Customer = customer
		comments = append(comments, comment)
	}

	err = rows.Err()
	return
}

func AddComment(articleId int64, comment Comment) error {
	_, err := database.Exec(`
		INSERT INTO comment
		(article, text, customer)
		VALUES
		(?, ?, ?)`, &articleId, &comment.Text, &comment.Customer.Id)
	return err
}
