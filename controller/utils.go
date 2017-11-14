package controller

import (
	"errors"
)

// IsFilled checks if all fields are filled, else it'll return an error.
func IsFilled(fields ...string) error {
	for _, field := range fields {
		if field == "" {
			return errors.New("all fields must be filled")
		}
	}

	return nil
}
