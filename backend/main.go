package main

import (
	"log"
	"os"
	"tranquil-pages/auth"
	"tranquil-pages/controllers"
	"tranquil-pages/database"
	"tranquil-pages/repository"
	"tranquil-pages/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func setupRoutes(db *database.Database) *gin.Engine {
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
	stateRepo := auth.NewOAuthStateRepository(db)
	tokenRepo := auth.NewTokenRepository(db)
	authService := auth.NewAuthService(auth.OAuthConfig, stateRepo, tokenRepo)
	authController := auth.NewAuthController(authService)

	// Setup router
	router := gin.Default()

	frontendURL, exists := os.LookupEnv("FRONTEND_URL")
	if !exists {
		log.Fatal("Failed to get FRONTEND_URL from environment")
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Setup public routes
	authController.SetupAuthRoutes(router)

	// Setup user api routes
	userApi := router.Group("/api")
	userApi.Use(auth.AuthMiddleware(authService))
	bookController.SetupBookRoutes(userApi)

	return router
}

func main() {
	// Initialize database
	db, err := database.GetDatabase()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	router := setupRoutes(db)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
