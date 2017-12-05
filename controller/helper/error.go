package helper

import (
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

const (
	DuplicateKeySqlError = 1062
)

func GetSqlErrorCode(err error) int {
	if sqlErr, ok := err.(*mysql.MySQLError); ok {
		return int(sqlErr.Number) // uint16 will always fit inside an int
	}

	return -1
}

func InternalServerError() error {
	return errors.New("Something internal went wrong!")
}

func InvalidFormDataError() error {
	return errors.New("Invalid form data")
}

func EmailAlreadyRegisteredError(email string) error {
	return fmt.Errorf("The email address '%s' is already registered", email)
}
