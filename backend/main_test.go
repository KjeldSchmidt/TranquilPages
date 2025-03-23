package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"tranquil-pages/auth"
	"tranquil-pages/database"

	"github.com/stretchr/testify/assert"
)

func TestBookRoutesProtected(t *testing.T) {
	// Initialize test database
	db, err := database.GetDatabase()
	if err != nil {
		t.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Get the actual router setup from main
	router := setupRoutes(db)

	// Test request without JWT
	t.Run("without JWT", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/books", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	// Test request with valid JWT
	t.Run("with valid JWT", func(t *testing.T) {
		testUser := &auth.GoogleUserInfo{
			ID:            "123",
			Email:         "test@example.com",
			VerifiedEmail: true,
		}
		token, _ := auth.GenerateToken(testUser)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/books", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
