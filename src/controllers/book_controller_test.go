package controllers

import (
	"betterreads/src/database"
	"betterreads/src/models"
	"betterreads/src/services"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getTestDependencies() *gin.Engine {
	router := gin.Default()

	db := database.GetTestDatabase()
	bookService := services.NewBookService(db)
	bookController := NewBookController(bookService)
	bookController.SetupBookRoutes(router)
	return router
}

func TestBookController_GivenEmptyDatabase_ReturnsNoBooks(t *testing.T) {
	router := getTestDependencies()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)
	router.ServeHTTP(w, req)

	var books []models.Book
	assert.Equal(t, http.StatusOK, w.Code)

	_ = json.Unmarshal(w.Body.Bytes(), &books)
	assert.Equal(t, 0, len(books))
}

func TestBookController_GivenBookIsCreated_ReturnsThatBook(t *testing.T) {
	// given
	router := getTestDependencies()
	book := models.Book{Title: "The Go Programming Language", Author: "Rob Pike", Comment: "sus", Rating: 2}

	w := httptest.NewRecorder()
	bookJson, _ := json.Marshal(book)
	req, _ := http.NewRequest("POST", "/books", bytes.NewReader(bookJson))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// when
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/books", nil)
	router.ServeHTTP(w, req)

	// then
	var books []models.Book
	assert.Equal(t, w.Code, http.StatusOK)

	_ = json.Unmarshal(w.Body.Bytes(), &books)
	assert.Equal(t, 1, len(books))
	assert.Equal(t, "The Go Programming Language", books[0].Title)
	assert.Equal(t, "Rob Pike", books[0].Author)
	assert.Equal(t, "sus", books[0].Comment)
	assert.Equal(t, 2, books[0].Rating)
}
