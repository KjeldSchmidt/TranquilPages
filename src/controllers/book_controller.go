package controllers

import (
	appErrors "betterreads/src/errors"
	"betterreads/src/models"
	"betterreads/src/services"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BookController struct {
	bookService *services.BookService
}

func NewBookController(bookService *services.BookService) *BookController {
	return &BookController{bookService: bookService}
}

func (bc *BookController) SetupBookRoutes(router *gin.Engine) {
	router.POST("/books", bc.CreateBook)
	router.GET("/books", bc.ListBooks)
	router.GET("/books/:id", bc.GetBook)
}

func (bc *BookController) CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := bc.bookService.CreateBook(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, book)
}

func (bc *BookController) ListBooks(c *gin.Context) {
	var books, err = bc.bookService.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, books)
}

func (bc *BookController) GetBook(c *gin.Context) {
	var book, err = bc.bookService.GetBookById(c.Param("id"))
	switch {
	case errors.Is(err, appErrors.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Book with id %s not found", c.Param("id"))})
	case errors.Is(err, appErrors.ErrDatabase):
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	case err == nil:
		c.JSON(http.StatusOK, book)
	}
}
