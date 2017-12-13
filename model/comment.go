package model

import "log"

type Comment struct {
	Id   int64
	Text string
	User int64
}

func GetCommentsByArticleId(id int64) (comments []Comment, err error) {
	rows, err := database.Query(`
		SELECT text FROM comments WHERE comments.article = (?)`, &id)
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

func AddComment(comment Comment) error {
	_, err := database.Exec(`
		INSERT INTO comments
		(?, ?)`, &comment.Text, &comment.User)
	return err
}
