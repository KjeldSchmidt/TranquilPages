package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestAuthMiddleware(t *testing.T) {
	// Setup test environment
	setupTestEnv(t)
	router := setupTestRouter()

	// Create mock repositories
	mockTokenRepo := NewMockTokenRepository()
	mockStateRepo := NewMockOAuthStateRepository()

	// Create test OAuth config
	config := &oauth2.Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://oauth2.googleapis.com/token",
		},
	}

	// Create test auth service
	authService := NewAuthService(config, mockStateRepo, mockTokenRepo)

	// Test user
	testUser := &GoogleUserInfo{
		ID:            "123",
		Email:         "test@test.com",
		VerifiedEmail: true,
	}

	// Generate valid token
	validToken, err := GenerateToken(testUser)
	assert.NoError(t, err)

	// Test handler
	router.GET("/test", AuthMiddleware(authService), func(c *gin.Context) {
		claims, exists := c.Get("claims")
		assert.True(t, exists)
		assert.NotNil(t, claims)
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "missing authorization header",
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "/test", nil)
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Authorization header is required"}`,
		},
		{
			name: "invalid header format",
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "InvalidFormat")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid authorization header format"}`,
		},
		{
			name: "invalid token",
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer invalid.token.here")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid or expired token"}`,
		},
		{
			name: "valid token",
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer "+validToken)
				return req
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"success"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, tt.setupRequest())

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
