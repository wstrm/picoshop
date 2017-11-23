SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

/* Create picoshopd schema */
CREATE SCHEMA IF NOT EXISTS picoshop DEFAULT CHARACTER SET utf8;
USE picoshop;

/* Create address table */
CREATE TABLE IF NOT EXISTS picoshop.address (
	PRIMARY KEY (id),

	id INT AUTO_INCREMENT,
	street VARCHAR(255) NOT NULL,
	care_of VARCHAR(255) NULL,
	zip_code INT(11) NOT NULL,
	country VARCHAR(255) NOT NULL);

/* Create user_has_address table */
CREATE TABLE IF NOT EXISTS picoshop.user_has_address (
	PRIMARY KEY (id),
	FOREIGN KEY (address) REFERENCES picoshop.address(id),

	id INT AUTO_INCREMENT,
	address INT NOT NULL);

/* Create user table */
CREATE TABLE IF NOT EXISTS picoshop.user (
	PRIMARY KEY (id),
	FOREIGN KEY (addresses)
		REFERENCES picoshop.user_has_address(id)
		ON DELETE CASCADE,

	UNIQUE INDEX uc_email (email ASC),

	id INT AUTO_INCREMENT,
	email VARCHAR(255) NOT NULL,
	hash VARCHAR(255) NOT NULL,
	name VARCHAR(255) NOT NULL,
	phone_number INT,
	create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	addresses INT NOT NULL);

/* Create comments table */
CREATE TABLE IF NOT EXISTS picoshop.comments (
	PRIMARY KEY (id),
	FOREIGN KEY (customer)
		REFERENCES picoshop.customer(id),

	id INT NOT NULL,
	rating DECIMAL(10, 0) NOT NULL, -- rating between 0-10 stars
	text VARCHAR(255),
	create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	customer INT NOT NULL);

/* Create article table */
CREATE TABLE IF NOT EXISTS picoshop.article (
	PRIMARY KEY (id),
	FOREIGN KEY (comments)
		REFERENCES picoshop.comments(id),

	id INT NOT NULL,
	name VARCHAR(255) NOT NULL,
	description VARCHAR(255) NOT NULL,
	price DECIMAL(11, 4) NOT NULL, -- decimal with 11 bits and 4 decimals
	image_url VARCHAR(255) NOT NULL,
	comments INT);

/* Create order_has_articles table */
CREATE TABLE IF NOT EXISTS picoshop.order_has_articles (
	PRIMARY KEY(id),
	FOREIGN KEY(article)
		REFERENCES picoshop.article(id),

	id INT NOT NULL,
	article INT NOT NULL
	);

/* Create order table */
CREATE TABLE IF NOT EXISTS picoshop.order (
	PRIMARY KEY (id),
	FOREIGN KEY (customer)
		REFERENCES picoshop.customer(id),
	FOREIGN KEY (address)
		REFERENCES picoshop.address(id),
	FOREIGN KEY (articles)
		REFERENCES picoshop.order_has_articles(id),

	id INT AUTO_INCREMENT,
	customer INT NOT NULL,
	address INT NOT NULL,
	status INT(3) NOT NULL, -- status, future proof with 3 long int
	articles INT NOT NULL,
	create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP);

/* Create user_has_orders table */
CREATE TABLE IF NOT EXISTS picoshop.customer_has_orders (
	PRIMARY KEY (id),
	FOREIGN KEY (`order`)
		REFERENCES picoshop.order(id),

	id INT AUTO_INCREMENT,
	`order` INT);

/* Create customer table */
CREATE TABLE IF NOT EXISTS picoshop.customer (
	PRIMARY KEY (id),
	FOREIGN KEY (user)
		REFERENCES picoshop.user(id)
		ON DELETE CASCADE,
	FOREIGN KEY (orders)
		REFERENCES picoshop.customer_has_orders(id),

	id INT AUTO_INCREMENT,
	password VARCHAR(255) NOT NULL,
	credit_card INT NOT NULL,
	user INT NOT NULL,
	orders INT);

/* Create admin table */
CREATE TABLE IF NOT EXISTS picoshop.admin (
	PRIMARY KEY (id),
	FOREIGN KEY (user)
		REFERENCES picoshop.user(id)
		ON DELETE CASCADE,

	id INT AUTO_INCREMENT NOT NULL,
	user INT NOT NULL);

/* Create warehouse table */
CREATE TABLE IF NOT EXISTS picoshop.warehouse (
	PRIMARY KEY (id),
	FOREIGN KEY (user)
		REFERENCES picoshop.user(id)
		ON DELETE CASCADE,

	id INT AUTO_INCREMENT NOT NULL,
	user INT NOT NULL);

SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
