package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestLogin(t *testing.T) {
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

	tests := []struct {
		name           string
		setupMock      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful login",
			setupMock: func() {
				// No setup needed for success case
			},
			expectedStatus: http.StatusTemporaryRedirect,
			expectedBody:   "",
		},
		{
			name: "service error",
			setupMock: func() {
				// Force service error by setting createFunc to return error
				mockStateRepo.createFunc = func(state *OAuthState) error {
					return assert.AnError
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Failed to generate redirect url for OAuth flow"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/auth/login", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestCallback(t *testing.T) {
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

	// Create test user info
	testUser := &GoogleUserInfo{
		ID:            "123",
		Email:         "test@test.com",
		VerifiedEmail: true,
	}

	// Create test server for userinfo endpoint
	userInfoServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/oauth2/v2/userinfo", r.URL.Path)
		json.NewEncoder(w).Encode(testUser)
	}))
	defer userInfoServer.Close()

	// Create test server for token endpoint
	tokenServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/token", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"access_token":  "test-access-token",
			"token_type":    "Bearer",
			"expires_in":    3600,
			"refresh_token": "test-refresh-token",
		})
	}))
	defer tokenServer.Close()

	// Update config to use test servers
	config.Endpoint.TokenURL = tokenServer.URL + "/token"
	config.Endpoint.AuthURL = tokenServer.URL + "/auth"

	// Override the user info URL in the auth service
	authService.userInfoURL = userInfoServer.URL + "/oauth2/v2/userinfo"

	tests := []struct {
		name           string
		setupMock      func()
		query          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "successful callback",
			setupMock: func() {
				// Store valid state
				mockStateRepo.Create(&OAuthState{
					State: "valid-state",
				})
			},
			query:          "?code=valid-code&state=valid-state",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"user":{"id":"123","email":"test@test.com","verified_email":true,"name":"","given_name":"","family_name":"","picture":"","locale":""},"token":"*"}`,
		},
		{
			name:           "missing code",
			setupMock:      func() {},
			query:          "?state=valid-state",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Code not found"}`,
		},
		{
			name:           "missing state",
			setupMock:      func() {},
			query:          "?code=valid-code",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"State parameter is required"}`,
		},
		{
			name: "invalid state",
			setupMock: func() {
				// Don't store any state
			},
			query:          "?code=valid-code&state=invalid-state",
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"OAuth state validation error: invalid or expired state"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/auth/callback"+tt.query, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedBody != "" {
				if tt.expectedBody == `{"user":{"id":"123","email":"test@test.com","verified_email":true,"name":"","given_name":"","family_name":"","picture":"","locale":""},"token":"*"}` {
					// For success case, just verify the structure and user info
					var response map[string]interface{}
					err := json.Unmarshal(w.Body.Bytes(), &response)
					assert.NoError(t, err)
					assert.Contains(t, response, "user")
					assert.Contains(t, response, "token")
					user := response["user"].(map[string]interface{})
					assert.Equal(t, "123", user["id"])
					assert.Equal(t, "test@test.com", user["email"])
					assert.True(t, user["verified_email"].(bool))
				} else {
					assert.JSONEq(t, tt.expectedBody, w.Body.String())
				}
			}
		})
	}
}

func TestLogout(t *testing.T) {
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
			name: "successful logout",
			setupMock: func() {
				// No setup needed for success case
			},
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("POST", "/auth/logout", nil)
				req.Header.Set("Authorization", "Bearer "+validToken)
				return req
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{}`,
		},
		{
			name:      "missing authorization header",
			setupMock: func() {},
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("POST", "/auth/logout", nil)
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   `{"error":"Authorization header is required"}`,
		},
		{
			name:      "invalid token format",
			setupMock: func() {},
			setupRequest: func() *http.Request {
				req, _ := http.NewRequest("POST", "/auth/logout", nil)
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
				req, _ := http.NewRequest("POST", "/auth/logout", nil)
				req.Header.Set("Authorization", "Bearer invalid.token.here")
				return req
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Failed to logout"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			w := httptest.NewRecorder()
			router.ServeHTTP(w, tt.setupRequest())

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, tt.expectedBody, w.Body.String())
		})
	}
}
