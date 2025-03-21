package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
)

func TestGetAuthURL(t *testing.T) {
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

	tests := []struct {
		name           string
		setupMock      func()
		expectedError  bool
		validateResult func(*testing.T, string, error)
	}{
		{
			name: "successful auth URL generation",
			setupMock: func() {
				// No special setup needed for success case
			},
			expectedError: false,
			validateResult: func(t *testing.T, url string, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, url)
				assert.Contains(t, url, "https://accounts.google.com/o/oauth2/auth")
				assert.Contains(t, url, "client_id=test-client-id")
				assert.Contains(t, url, "redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback")

				// Verify state was stored
				states := mockStateRepo.states
				assert.Len(t, states, 1)
				for _, state := range states {
					assert.NotEmpty(t, state.State)
					assert.NotZero(t, state.CreatedAt)
					assert.NotZero(t, state.ExpiresAt)
				}
			},
		},
		{
			name: "state repository error",
			setupMock: func() {
				// Simulate repository error
				mockStateRepo.createFunc = func(state *OAuthState) error {
					return assert.AnError
				}
			},
			expectedError: true,
			validateResult: func(t *testing.T, url string, err error) {
				assert.Error(t, err)
				assert.Empty(t, url)
				assert.IsType(t, &AuthURLGenerationError{}, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			url, err := authService.GetAuthURL()
			tt.validateResult(t, url, err)
		})
	}
}

func TestHandleCallback(t *testing.T) {
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

	// Create test user info
	testUser := &GoogleUserInfo{
		ID:            "123",
		Email:         "test@test.com",
		VerifiedEmail: true,
		Name:          "Test User",
		GivenName:     "Test",
		FamilyName:    "User",
		Picture:       "https://example.com/picture.jpg",
		Locale:        "en",
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
		code           string
		state          string
		expectedError  bool
		validateResult func(*testing.T, *GoogleUserInfo, error)
	}{
		{
			name: "successful callback",
			setupMock: func() {
				// Store valid state
				now := time.Now()
				mockStateRepo.Create(&OAuthState{
					State:     "valid-state",
					CreatedAt: primitive.NewDateTimeFromTime(now),
					ExpiresAt: primitive.NewDateTimeFromTime(now.Add(15 * time.Minute)),
				})
			},
			code:          "valid-code",
			state:         "valid-state",
			expectedError: false,
			validateResult: func(t *testing.T, user *GoogleUserInfo, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, testUser.ID, user.ID)
				assert.Equal(t, testUser.Email, user.Email)
				assert.Equal(t, testUser.VerifiedEmail, user.VerifiedEmail)
				assert.Equal(t, testUser.Name, user.Name)
				assert.Equal(t, testUser.GivenName, user.GivenName)
				assert.Equal(t, testUser.FamilyName, user.FamilyName)
				assert.Equal(t, testUser.Picture, user.Picture)
				assert.Equal(t, testUser.Locale, user.Locale)
			},
		},
		{
			name: "invalid state",
			setupMock: func() {
				// Don't store any state
			},
			code:          "valid-code",
			state:         "invalid-state",
			expectedError: true,
			validateResult: func(t *testing.T, user *GoogleUserInfo, err error) {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.IsType(t, &StateValidationError{}, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			user, err := authService.HandleCallback(tt.code, tt.state)
			tt.validateResult(t, user, err)
		})
	}
}
