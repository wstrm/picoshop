package controller

import (
	"github.com/go-sql-driver/mysql"
)

const (
	DuplicateKeySqlError = 1062
)

func getSqlErrorCode(err error) int {
	if sqlErr, ok := err.(*mysql.MySQLError); ok {
		return int(sqlErr.Number) // uint16 will always fit inside an int
	}

	return -1
}
