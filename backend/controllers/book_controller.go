package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"tranquil-pages/auth"
	appErrors "tranquil-pages/errors"
	"tranquil-pages/models"
	"tranquil-pages/services"

	"github.com/gin-gonic/gin"
)

type BookController struct {
	bookService *services.BookService
}

func NewBookController(bookService *services.BookService) *BookController {
	return &BookController{bookService: bookService}
}

func (bc *BookController) SetupBookRoutes(router *gin.RouterGroup) {
	router.POST("/books", bc.CreateBook)
	router.GET("/books", bc.ListBooks)
	router.GET("/books/:id", bc.GetBook)
	router.DELETE("/books/:id", bc.DeleteBook)
}

func (bc *BookController) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, appErrors.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Book with id %s not found", c.Param("id"))})
	case errors.Is(err, appErrors.ErrDatabase):
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func (bc *BookController) CreateBook(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book.UserID = claims.UserID
	err := bc.bookService.CreateBook(&book)
	if err != nil {
		bc.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, book)
}

func (bc *BookController) ListBooks(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	books, err := bc.bookService.GetBooksByUserID(claims.UserID)
	if err != nil {
		bc.handleError(c, err)
		return
	}

	if books == nil {
		books = []models.Book{} // Ensure an empty slice instead of nil
	}

	c.JSON(http.StatusOK, books)
}

func (bc *BookController) GetBook(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	book, err := bc.bookService.GetBook(c.Param("id"), claims.UserID)
	if err != nil {
		bc.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, book)
}

func (bc *BookController) DeleteBook(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	err := bc.bookService.DeleteBook(c.Param("id"), claims.UserID)
	if err != nil {
		bc.handleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
