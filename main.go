package main

import (
	"betterreads/src/controllers"
	"betterreads/src/database"
	"betterreads/src/services"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	db, err := database.GetDbHandler()
	if err != nil {
		log.Fatalf("Error setting up database: %v", err)
	}

	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	bookService := services.NewBookService(db)
	bookController := controllers.NewBookController(bookService)

	bookController.SetupBookRoutes(router)

	err = router.SetTrustedProxies(nil)
	if err != nil {
		panic("We... failed at not trusting any proxies...? I guess?")
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
