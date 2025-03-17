package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"tranquil-pages/src/database"
	"tranquil-pages/src/models"
	"tranquil-pages/src/services"
	"tranquil-pages/src/test_utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getTestDependencies() (*gin.Engine, *database.TestDatabase) {
	router := gin.Default()

	testDB, err := database.NewTestDatabase()
	if err != nil {
		panic(err)
	}

	bookService := services.NewBookService(testDB.GetDatabase())
	bookController := NewBookController(bookService)
	bookController.SetupBookRoutes(router)
	return router, testDB
}

func makeRandomBook() *models.Book {
	author := "Author " + test_utils.RandomString(12)
	title := "Title " + test_utils.RandomString(20)
	comment := "Comment " + test_utils.RandomString(20)
	rating := rand.Intn(6)
	book := models.Book{Author: author, Title: title, Comment: comment, Rating: rating}
	return &book
}

func createBookViaApi(router *gin.Engine, book *models.Book) *models.Book {
	w := httptest.NewRecorder()
	bookJson, _ := json.Marshal(book)
	req, _ := http.NewRequest("POST", "/books", bytes.NewReader(bookJson))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	var resultBook *models.Book
	_ = json.Unmarshal(w.Body.Bytes(), &resultBook)
	return resultBook
}

func TestBookController_GivenEmptyDatabase_ReturnsNoBooks(t *testing.T) {
	// Given
	router, testDB := getTestDependencies()
	defer testDB.Close()

	// When
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)
	router.ServeHTTP(w, req)

	// Then
	var books []models.Book
	assert.Equal(t, http.StatusOK, w.Code)

	assert.Equal(t, "[]", w.Body.String()) // We want to return an empty list, not an empty body or null object

	_ = json.Unmarshal(w.Body.Bytes(), &books)
	assert.Equal(t, 0, len(books))
}

func TestBookController_GivenBookIsCreated_ReturnsThatBook(t *testing.T) {
	// given
	router, testDB := getTestDependencies()
	defer testDB.Close()
	book := makeRandomBook()

	createBookViaApi(router, book)

	// when
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)
	router.ServeHTTP(w, req)

	// then
	var books []models.Book
	assert.Equal(t, w.Code, http.StatusOK)

	_ = json.Unmarshal(w.Body.Bytes(), &books)
	assert.Equal(t, 1, len(books))
	assert.True(t, models.CompareBooks(book, &books[0]))
}

func TestBookController_GivenBookIsCreated_CanFetchThatBookById(t *testing.T) {
	// given
	router, testDB := getTestDependencies()
	defer testDB.Close()
	transientBook := makeRandomBook()
	expectedBook := createBookViaApi(router, transientBook)

	// when
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/books/%s", expectedBook.ID.Hex()), nil)
	router.ServeHTTP(w, req)

	// then
	var actualBook *models.Book
	assert.Equal(t, w.Code, http.StatusOK)

	_ = json.Unmarshal(w.Body.Bytes(), &actualBook)
	assert.True(t, models.CompareBooks(expectedBook, actualBook))
}

func TestBookController_GivenManyBooksAreCreated_CanFetchSpecificBookById(t *testing.T) {
	// given
	router, testDB := getTestDependencies()
	defer testDB.Close()
	bookCount := 5

	var expectedBook *models.Book
	for i := 0; i < bookCount; i++ {
		transientBook := makeRandomBook()
		expectedBook = createBookViaApi(router, transientBook)
	}

	// when
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/books/%s", expectedBook.ID.Hex()), nil)
	router.ServeHTTP(w, req)

	// then
	var actualBook *models.Book
	assert.Equal(t, w.Code, http.StatusOK)

	_ = json.Unmarshal(w.Body.Bytes(), &actualBook)
	assert.True(t, models.CompareBooks(expectedBook, actualBook))
}

func TestBookController_GivenManyBooksAreCreated_ReturnsAllBooks(t *testing.T) {
	// given
	router, testDB := getTestDependencies()
	defer testDB.Close()
	bookCount := 5

	expectedBooks := make(map[primitive.ObjectID]*models.Book)
	for i := 0; i < bookCount; i++ {
		transientBook := makeRandomBook()
		book := createBookViaApi(router, transientBook)
		expectedBooks[book.ID] = book
	}

	// when
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)
	router.ServeHTTP(w, req)

	// then
	var actualBooks []models.Book
	assert.Equal(t, w.Code, http.StatusOK)
	_ = json.Unmarshal(w.Body.Bytes(), &actualBooks)

	for _, actualBook := range actualBooks {
		expectedBook, _ := expectedBooks[actualBook.ID]
		assert.True(t, models.CompareBooks(expectedBook, &actualBook))
	}
	assert.Equal(t, len(expectedBooks), len(actualBooks))
}

func TestBookController_CanDeleteExistingBook(t *testing.T) {
	// given
	router, testDB := getTestDependencies()
	defer testDB.Close()
	book := makeRandomBook()
	createdBook := createBookViaApi(router, book)

	// when
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/books/%s", createdBook.ID.Hex()), nil)
	router.ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusNoContent, w.Code)

	// and then
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/books", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, "[]", w.Body.String())
}

func TestBookController_CanDeleteNonExistentBook(t *testing.T) {
	// given
	router, testDB := getTestDependencies()
	defer testDB.Close()
	nonExistentID := primitive.NewObjectID()

	// when
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/books/%s", nonExistentID.Hex()), nil)
	router.ServeHTTP(w, req)

	// then
	assert.Equal(t, http.StatusNoContent, w.Code)
}
