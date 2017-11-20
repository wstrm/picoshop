-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='TRADITIONAL,ALLOW_INVALID_DATES';

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `mydb` DEFAULT CHARACTER SET utf8 ;
USE `mydb` ;

-- -----------------------------------------------------
-- Table `mydb`.`Adress`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`Adress` (
  `id` INT(11) NOT NULL,
  `street_adress` VARCHAR(64) NOT NULL,
  `care_of` VARCHAR(64) NULL,
  `zip_code` INT(11) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `country` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mydb`.`User_has_Adress`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`User_has_Adress` (
  `Adress_id` INT(11) NOT NULL,
  `Adress_id1` INT(11) NOT NULL,
  PRIMARY KEY (`Adress_id`),
  INDEX `fk_User_has_Adress_Adress1_idx` (`Adress_id1` ASC),
  CONSTRAINT `fk_User_has_Adress_Adress1`
    FOREIGN KEY (`Adress_id1`)
    REFERENCES `mydb`.`Adress` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mydb`.`User`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`User` (
  `id` INT(11) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `hash` VARCHAR(32) NOT NULL,
  `first_name` VARCHAR(32) NOT NULL,
  `last_name` VARCHAR(32) NOT NULL,
  `phone_nr` INT(11) NULL,
  `adress_adress_id` INT(11) NULL,
  `name` VARCHAR(45) NOT NULL,
  `date` DATETIME NOT NULL,
  `User_has_Adress_Adress_id` INT(11) NOT NULL,
  PRIMARY KEY (`id`, `adress_adress_id`),
  UNIQUE INDEX `email_UNIQUE` (`email` ASC),
  INDEX `fk_User_User_has_Adress1_idx` (`User_has_Adress_Adress_id` ASC),
  CONSTRAINT `fk_User_User_has_Adress1`
    FOREIGN KEY (`User_has_Adress_Adress_id`)
    REFERENCES `mydb`.`User_has_Adress` (`Adress_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB
COMMENT = 'on remove - cascade\non update - cascade';


-- -----------------------------------------------------
-- Table `mydb`.`Order`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`Order` (
  `id` INT(11) NOT NULL,
  `Customer_id` VARCHAR(32) NULL,
  `Customer_User_id` INT(11) NULL,
  `Customer_User_Adress_id` INT(11) NULL,
  `status` VARCHAR(45) NOT NULL,
  `date` DATETIME NOT NULL,
  PRIMARY KEY (`id`, `Customer_id`, `Customer_User_id`, `Customer_User_Adress_id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mydb`.`Customer`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`Customer` (
  `id` VARCHAR(32) NOT NULL,
  `password` VARCHAR(32) NOT NULL,
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `credit_card` INT(11) NOT NULL,
  `User_id` INT(11) NOT NULL,
  `User_Adress_id` INT(11) NULL,
  `Order_id` INT(11) NOT NULL,
  `Order_Customer_id` VARCHAR(32) NOT NULL,
  `Order_Customer_User_id` INT(11) NOT NULL,
  `Order_Customer_User_Adress_id` INT(11) NOT NULL,
  PRIMARY KEY (`id`, `User_id`, `User_Adress_id`),
  INDEX `fk_customer_user1_idx` (`User_id` ASC, `User_Adress_id` ASC),
  INDEX `fk_Customer_Order1_idx` (`Order_id` ASC, `Order_Customer_id` ASC, `Order_Customer_User_id` ASC, `Order_Customer_User_Adress_id` ASC),
  CONSTRAINT `fk_customer_user1`
    FOREIGN KEY (`User_id` , `User_Adress_id`)
    REFERENCES `mydb`.`User` (`id` , `adress_adress_id`)
    ON DELETE CASCADE
    ON UPDATE CASCADE,
  CONSTRAINT `fk_Customer_Order1`
    FOREIGN KEY (`Order_id` , `Order_Customer_id` , `Order_Customer_User_id` , `Order_Customer_User_Adress_id`)
    REFERENCES `mydb`.`Order` (`id` , `Customer_id` , `Customer_User_id` , `Customer_User_Adress_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mydb`.`Admin`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`Admin` (
  `id` VARCHAR(16) NOT NULL,
  `User_id` INT(11) NOT NULL,
  `User_adress_adress_id` INT(11) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_Admin_User1_idx` (`User_id` ASC, `User_adress_adress_id` ASC),
  CONSTRAINT `fk_Admin_User1`
    FOREIGN KEY (`User_id` , `User_adress_adress_id`)
    REFERENCES `mydb`.`User` (`id` , `adress_adress_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mydb`.`Warehouse`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`Warehouse` (
  `id` VARCHAR(16) NOT NULL,
  `User_id` INT(11) NOT NULL,
  `User_adress_adress_id` INT(11) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_Warehouse_User1_idx` (`User_id` ASC, `User_adress_adress_id` ASC),
  CONSTRAINT `fk_Warehouse_User1`
    FOREIGN KEY (`User_id` , `User_adress_adress_id`)
    REFERENCES `mydb`.`User` (`id` , `adress_adress_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mydb`.`Articles`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`Articles` (
  `id` INT NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `description` VARCHAR(255) NULL,
  `price` DECIMAL(11,4) NOT NULL,
  `image_url` VARCHAR(255) NULL,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mydb`.`Comments`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`Comments` (
  `id` INT NOT NULL,
  `rating` DECIMAL(10,0) NOT NULL,
  `text` VARCHAR(255) NULL,
  `date` DATETIME NOT NULL,
  `Customer_id` VARCHAR(32) NOT NULL,
  `Customer_User_id` INT(11) NOT NULL,
  `Customer_User_Adress_id` INT(11) NOT NULL,
  `Articles_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_Comments_Customer1_idx` (`Customer_id` ASC, `Customer_User_id` ASC, `Customer_User_Adress_id` ASC),
  INDEX `fk_Comments_Articles1_idx` (`Articles_id` ASC),
  CONSTRAINT `fk_Comments_Customer1`
    FOREIGN KEY (`Customer_id` , `Customer_User_id` , `Customer_User_Adress_id`)
    REFERENCES `mydb`.`Customer` (`id` , `User_id` , `User_Adress_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_Comments_Articles1`
    FOREIGN KEY (`Articles_id`)
    REFERENCES `mydb`.`Articles` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `mydb`.`Order_has_Articles`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `mydb`.`Order_has_Articles` (
  `id` INT NOT NULL,
  `Order_id` INT(11) NOT NULL,
  `Order_Customer_id` VARCHAR(32) NOT NULL,
  `Order_Customer_User_id` INT(11) NOT NULL,
  `Order_Customer_User_Adress_id` INT(11) NULL,
  `Articles_id` INT NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_Order_has_Articles_Order1_idx` (`Order_id` ASC, `Order_Customer_id` ASC, `Order_Customer_User_id` ASC, `Order_Customer_User_Adress_id` ASC),
  INDEX `fk_Order_has_Articles_Articles1_idx` (`Articles_id` ASC),
  CONSTRAINT `fk_Order_has_Articles_Order1`
    FOREIGN KEY (`Order_id` , `Order_Customer_id` , `Order_Customer_User_id` , `Order_Customer_User_Adress_id`)
    REFERENCES `mydb`.`Order` (`id` , `Customer_id` , `Customer_User_id` , `Customer_User_Adress_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_Order_has_Articles_Articles1`
    FOREIGN KEY (`Articles_id`)
    REFERENCES `mydb`.`Articles` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
