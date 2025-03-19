package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound      = errors.New("Record not found")
	ErrDatabase      = errors.New("Database Error")
	ErrInvalidID     = errors.New("Invalid ID format")
	ErrInvalidRating = errors.New("Rating must be between 0 and 5")
	ErrDuplicateBook = errors.New("A book with this title already exists")
	ErrConnection    = errors.New("Failed to connect to database")
)

func ErrEnvNotSet(varName string) error {
	return fmt.Errorf("environment variable %s not set", varName)
}
