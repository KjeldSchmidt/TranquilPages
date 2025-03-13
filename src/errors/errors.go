package errors

import "errors"

var (
	ErrNotFound = errors.New("Record not found")
	ErrDatabase = errors.New("Database Error")
)
