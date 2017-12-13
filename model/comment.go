package model

import "log"

type Comment struct {
	Id       int64
	Text     string
	Customer int64
}

func GetCommentsByArticleId(id int64) (comments []Comment, err error) {
	rows, err := database.Query(`
		SELECT text FROM comment WHERE comment.article = (?)`, &id)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		comment := Comment{}

		err = rows.Scan(
			&comment.Text)
		if err != nil {
			log.Panicln(err)
		}

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
		(?, ?, ?)`, &articleId, &comment.Text, &comment.Customer)
	return err
}
