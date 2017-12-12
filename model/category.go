package model

import "log"

type Category struct {
	Name          string
	Subcategories []Subcategory
}

type Subcategory struct {
	Name     string
	Category string
	Articles []Article
}

func NewCategory(name string) Category {
	return Category{
		Name: name,
	}
}

func NewSubcategory(name string) Subcategory {
	return Subcategory{
		Name: name,
	}
}

func GetAllCategories() (categories []Category, err error) {
	rows, err := database.Query(`
		SELECT category, subcategory
		FROM category_has_subcategories
		ORDER BY category
		`)
	if err != nil {
		return
	}

	defer rows.Close()

	c := make(map[string][]Subcategory)
	var k, v string
	var s []string // save order of map

	for rows.Next() {
		err = rows.Scan(&k, &v)
		if err != nil {
			log.Panicln(err)
		}

		// do not add key to order slice if already exists
		if _, exists := c[k]; !exists {
			s = append(s, k)
		}

		c[k] = append(c[k], NewSubcategory(v))
	}
	err = rows.Err()
	if err != nil {
		return
	}

	for _, name := range s { // range over ordered slice
		categories = append(categories, Category{
			Name:          name,
			Subcategories: c[name],
		})
	}

	return
}

func GetSubcategoriesFromCategory(category Category) (subcategories []Subcategory, err error) {
	rows, err := database.Query(`
		SELECT (subcategory.name, subcategory.articles)
		FROM category WHERE category.id=?

		INNER JOIN subcategories
		ON category.subcategories = category_has_subcategories.id

		INNER JOIN subcategory
		ON category_has_subcategories.subcategory = subcategory.id

		ORDER BY category.name
	`, category.Name)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		subcategory := Subcategory{}

		err = rows.Scan(
			&subcategory.Name, &subcategory.Articles)
		if err != nil {
			log.Panicln(err)
		}

		subcategories = append(subcategories, subcategory)
	}

	err = rows.Err()
	return
}

func putCategory(category Category) (Category, error) {
	_, err := database.Exec(`
		INSERT IGNORE INTO category
		(name)
		VALUES
		(TRIM(LOWER(?)))
	`, category.Name)

	return category, err
}

func putSubcategory(subcategory Subcategory) (Subcategory, error) {
	_, err := database.Exec(`
		INSERT IGNORE INTO subcategory
		(name, category)
		VALUES
		(TRIM(LOWER(?)), TRIM(LOWER(?)))
	`, subcategory.Name, subcategory.Category)

	return subcategory, err
}

func addSubcategoryToCategory(category Category, subcategory Subcategory) error {
	_, err := database.Exec(`
		INSERT IGNORE INTO category_has_subcategories
		(category, subcategory)
		VALUES
		(TRIM(LOWER(?)), TRIM(LOWER(?)))
	`, category.Name, subcategory.Name)

	return err
}

func ensureSubcategoryWithCategory(category Category, subcategory Subcategory) {
	cat, err := putCategory(category)
	if err != nil {
		log.Panicln(err)
	}

	subcategory.Category = cat.Name
	subcat, err := putSubcategory(subcategory)
	if err != nil {
		log.Panicln(err)
	}

	err = addSubcategoryToCategory(cat, subcat)
	if err != nil {
		log.Panicln(err)
	}
}
