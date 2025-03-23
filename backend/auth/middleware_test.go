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

	// Create test controller
	controller := NewAuthController(authService)

	// Setup test router
	router := setupTestRouter()
	controller.SetupAuthRoutes(router)

	// Add test endpoint
	router.GET("/test", AuthMiddleware(authService), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	// Generate valid token for testing
	testUser := &GoogleUserInfo{
		ID:            "123",
		Email:         "test@test.com",
		VerifiedEmail: true,
	}
	validToken, err := GenerateToken(testUser)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		setupMock      func()
		setupRequest   func() *http.Request
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "missing authorization header",
			setupMock: func() {},
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "/test", nil)
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"no valid authentication token found"}`,
		},
		{
			name:      "invalid header format",
			setupMock: func() {},
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "InvalidFormat")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid authorization header format"}`,
		},
		{
			name:      "invalid token",
			setupMock: func() {},
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer invalid.token.here")
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Invalid or expired token"}`,
		},
		{
			name: "revoked token",
			setupMock: func() {
				mockTokenRepo.Blacklist(validToken)
			},
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer "+validToken)
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Token has been revoked"}`,
		},
		{
			name:      "valid token",
			setupMock: func() {},
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "/test", nil)
				req.Header.Set("Authorization", "Bearer "+validToken)
				return req
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"success"}`,
		},
		{
			name:      "valid token from cookie",
			setupMock: func() {},
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("GET", "/test", nil)
				req.AddCookie(&http.Cookie{
					Name:     "token",
					Value:    validToken,
					HttpOnly: true,
					Secure:   true,
				})
				return req
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"success"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mock repository state
			mockTokenRepo.Reset()
			tt.setupMock()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, tt.setupRequest())

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
