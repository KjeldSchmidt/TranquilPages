package main

import (
	"log"
	"tranquil-pages/auth"
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

	// Initialize OAuth
	if err := auth.InitOAuthConfig(); err != nil {
		log.Fatal("Failed to initialize OAuth config:", err)
	}
	authService := auth.NewAuthService(auth.OAuthConfig)
	authController := auth.NewAuthController(authService)

	// Setup router
	router := gin.Default()

	// Setup public routes
	authController.SetupAuthRoutes(router)

	// Setup user api routes
	user_api := router.Group("/api")
	bookController.SetupBookRoutes(user_api)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
