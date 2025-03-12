package database

import (
	"betterreads/src/models"
	"gorm.io/gorm/utils/tests"
	"testing"
)

func TestDatabaseWrite(t *testing.T) {
	dbHandler := GetTestDatabase()

	book := models.Book{Title: "The Go Programming Language", Author: "Rob Pike", Comment: "sus", Rating: 2}
	dbHandler.Create(&book)

	retrievedBook := models.Book{}
	dbHandler.First(&retrievedBook)
	tests.AssertEqual(t, retrievedBook.Title, "The Go Programming Language")
}
