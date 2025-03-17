package main

import (
	"log"
	"tranquil-pages/src/controllers"
	"tranquil-pages/src/database"
	"tranquil-pages/src/repository"
	"tranquil-pages/src/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db, err := database.GetDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repository
	bookRepo := repository.NewBookRepository(db)

	// Initialize service
	bookService := services.NewBookService(bookRepo)

	// Initialize controller
	bookController := controllers.NewBookController(bookService)

	// Setup router
	router := gin.Default()

	// Book routes
	router.POST("/books", bookController.CreateBook)
	router.GET("/books", bookController.ListBooks)
	router.GET("/books/:id", bookController.GetBook)
	router.DELETE("/books/:id", bookController.DeleteBook)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
