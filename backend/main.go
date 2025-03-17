package main

import (
	"log"
	"tranquil-pages/controllers"
	"tranquil-pages/database"
	"tranquil-pages/repository"
	"tranquil-pages/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	db, err := database.GetDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	bookRepo := repository.NewBookRepository(db)

	// Initialize services
	bookService := services.NewBookService(bookRepo)

	// Initialize controllers
	bookController := controllers.NewBookController(bookService)

	// Setup router
	router := gin.Default()

	// Setup routes
	bookController.SetupBookRoutes(router)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
