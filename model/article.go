package model

import "log"

type Article struct {
	Id          int64
	Name        string
	Description string
	Price       float64
	ImageName   string
	Category    string
	Subcategory string
}

func NewArticle(name, description string, price float64, imageName string, category string, subcategory string) Article {
	return Article{
		Name:        name,
		Description: description,
		Price:       price,
		ImageName:   imageName,
		Category:    category,
		Subcategory: subcategory,
	}
}

func SearchForArticles(query string) (articles []Article, err error) {
	rows, err := database.Query(`
		SELECT id, name, description, price, image_name, category, subcategory
		FROM article WHERE name LIKE ?`, "%"+query+"%")
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		article := Article{}

		err = rows.Scan(
			&article.Id, &article.Name, &article.Description,
			&article.Price, &article.ImageName, &article.Category, &article.Subcategory)
		if err != nil {
			log.Panicln(err)
		}

		articles = append(articles, article)
	}

	err = rows.Err()
	return
}

func PutArticle(article Article) (Article, error) {
	ensureSubcategoryWithCategory(NewCategory(article.Category), NewSubcategory(article.Subcategory))

	result, err := database.Exec(`
		INSERT INTO article
		(name, description, price, image_name, category, subcategory)
		VALUES
		(?, ?, ?, ?, ?, ?)
	`, &article.Name, &article.Description, &article.Price, &article.ImageName, &article.Category, &article.Subcategory)
	if err != nil {
		return Article{}, err
	}

	article.Id, err = result.LastInsertId()

	_, err = database.Exec(`
		INSERT INTO subcategory_has_articles
		(subcategory, article)
		VALUES
		(?, ?)
	`, &article.Subcategory, &article.Id)

	return article, err
}

func GetArticlesFromSubcategory(subcategory Subcategory) (articles []Article, err error) {
	rows, err := database.Query(`
		SELECT (article.id, article.name, article.description, article.price, article.image_name)
		FROM subcategory WHERE subcategory.id=?

		INNER JOIN articles
		ON subcategory.articles = subcategory_has_articles.id

		INNER JOIN article
		ON subcategory_has_articles.article = article.id
	`, subcategory.Name)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		article := Article{}

		err = rows.Scan(
			&article.Id, &article.Name, &article.Description, &article.Price, &article.ImageName)
		if err != nil {
			log.Panicln(err)
		}

		articles = append(articles, article)
	}

	err = rows.Err()
	return
}

func GetArticleHighlights(n uint) (articles []Article, err error) {
	rows, err := database.Query(`
		SELECT id, name, description, price, image_name, category, subcategory
		FROM article
		WHERE rand() <= 0.3
		LIMIT ?
	`, n)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		article := Article{}

		err = rows.Scan(
			&article.Id, &article.Name, &article.Description, &article.Price, &article.ImageName, &article.Category, &article.Subcategory)
		if err != nil {
			log.Panicln(err)
		}

		articles = append(articles, article)
	}

	err = rows.Err()
	return
}

func GetArticleById(id int64) (article Article, err error) {
	err = database.QueryRow(`
		SELECT id, name, description, price, image_name, category, subcategory
		FROM article
		WHERE id=?
	`, id).Scan(&article.Id, &article.Name, &article.Description, &article.Price, &article.ImageName, &article.Category, &article.Subcategory)

	return
}

func GetArticlesByCategory(category string) (articles []Article, err error) {
	rows, err := database.Query(`
		SELECT a.id, a.name, a.description, a.price, a.image_name, a.category, a.subcategory
		FROM category_has_subcategories c
		INNER JOIN subcategory_has_articles s
		ON c.subcategory=s.subcategory
		INNER JOIN article a
		ON s.article=a.id
		WHERE a.category=?
		ORDER BY a.name
	`, category)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		article := Article{}

		err = rows.Scan(
			&article.Id, &article.Name, &article.Description, &article.Price, &article.ImageName, &article.Category, &article.Subcategory)
		if err != nil {
			log.Panicln(err)
		}

		articles = append(articles, article)
	}

	err = rows.Err()
	return
}

func GetArticlesBySubcategory(subcategory string) (articles []Article, err error) {
	rows, err := database.Query(`
		SELECT a.id, a.name, a.description, a.price, a.image_name, a.category, a.subcategory
		FROM subcategory_has_articles s
		INNER JOIN article a
		ON s.article=a.id
		WHERE s.subcategory=?
		ORDER BY a.name
	`, subcategory)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		article := Article{}

		err = rows.Scan(
			&article.Id, &article.Name, &article.Description, &article.Price, &article.ImageName, &article.Category, &article.Subcategory)
		if err != nil {
			log.Panicln(err)
		}

		articles = append(articles, article)
	}

	err = rows.Err()
	return
}
