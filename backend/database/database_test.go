package database

import (
	"context"
	"testing"
	"time"
	"tranquil-pages/models"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDatabaseWrite(t *testing.T) {
	dbHandler, err := NewTestDatabase()
	if err != nil {
		t.Fatal(err)
	}
	defer dbHandler.Close()

	// Create test book
	book := models.Book{
		Title:     "The Go Programming Language",
		Author:    "Rob Pike",
		Comment:   "sus",
		Rating:    2,
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	// Insert book
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = dbHandler.GetCollection("books").InsertOne(ctx, book)
	if err != nil {
		t.Fatal(err)
	}

	// Retrieve book
	var retrievedBook models.Book
	err = dbHandler.GetCollection("books").FindOne(ctx, bson.M{"title": book.Title}).Decode(&retrievedBook)
	if err != nil {
		t.Fatal(err)
	}

	// Assert values
	assert.Equal(t, book.Title, retrievedBook.Title)
	assert.Equal(t, book.Author, retrievedBook.Author)
	assert.Equal(t, book.Comment, retrievedBook.Comment)
	assert.Equal(t, book.Rating, retrievedBook.Rating)
}
