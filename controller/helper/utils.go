package helper

import "fmt"

// IsFilled checks if all fields are filled, else it'll return an error.
func IsFilled(fields ...string) error {
	for id, field := range fields {
		if field == "" {
			return fmt.Errorf("field number: %d is not filled", id)
		}
	}

	return nil
}
