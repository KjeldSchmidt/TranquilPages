package auth

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setupTestEnv(t *testing.T) {
	// Set a test JWT secret
	secret := "dGVzdF9zZWNyZXRfZm9yX2p3dF90ZXN0aW5nXzMyYnl0ZXM=" // base64 encoded 32-byte test secret
	os.Setenv("JWT_SECRET", secret)
}

// generateRandomString creates a random string of specified length
func generateRandomString(length int) string {
	b := make([]byte, length)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:length]
}

func TestGenerateToken(t *testing.T) {
	setupTestEnv(t)

	tests := []struct {
		name     string
		user     *GoogleUserInfo
		wantErr  bool
		validate func(*testing.T, string, error)
	}{
		{
			name: "valid user",
			user: &GoogleUserInfo{
				ID:            "123",
				Email:         "test@test.com",
				VerifiedEmail: true,
			},
			wantErr: false,
			validate: func(t *testing.T, token string, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// Verify token can be validated
				claims, err := ValidateToken(token)
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, "123", claims.UserID)
				assert.Equal(t, "test@test.com", claims.Email)
				assert.True(t, claims.Verified)
			},
		},
		{
			name: "missing JWT secret",
			user: &GoogleUserInfo{
				ID:            "123",
				Email:         "test@test.com",
				VerifiedEmail: true,
			},
			wantErr: true,
			validate: func(t *testing.T, token string, err error) {
				assert.Error(t, err)
				assert.Empty(t, token)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "missing JWT secret" {
				os.Unsetenv("JWT_SECRET")
			} else {
				setupTestEnv(t)
			}

			token, err := GenerateToken(tt.user)
			tt.validate(t, token, err)
		})
	}
}

func TestValidateToken(t *testing.T) {
	setupTestEnv(t)

	staticUser := &GoogleUserInfo{
		ID:            "123",
		Email:         "test@test.com",
		VerifiedEmail: true,
	}
	validStaticToken, err := GenerateToken(staticUser)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		token    string
		wantErr  bool
		validate func(*testing.T, *Claims, error)
	}{
		{
			name:    "valid token",
			token:   validStaticToken,
			wantErr: false,
			validate: func(t *testing.T, claims *Claims, err error) {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
				assert.Equal(t, staticUser.ID, claims.UserID)
				assert.Equal(t, staticUser.Email, claims.Email)
				assert.Equal(t, staticUser.VerifiedEmail, claims.Verified)
			},
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
			validate: func(t *testing.T, claims *Claims, err error) {
				assert.Error(t, err)
				assert.Nil(t, claims)
			},
		},
		{
			name:    "invalid token format",
			token:   "invalid.token.format",
			wantErr: true,
			validate: func(t *testing.T, claims *Claims, err error) {
				assert.Error(t, err)
				assert.Nil(t, claims)
			},
		},
		{
			name:    "expired token",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidGVzdCIsImVtYWlsIjoidGVzdEBleGFtcGxlLmNvbSIsInZlcmlmaWVkIjp0cnVlLCJleHAiOjE2MTYyMjQ4MDAsImlhdCI6MTYxNjIxMTIwMH0.4Adcj3UFYzPUVaVF43FmMze0Qp6ZQwWqZqZQwWqZqZQ",
			wantErr: true,
			validate: func(t *testing.T, claims *Claims, err error) {
				assert.Error(t, err)
				assert.Nil(t, claims)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token)
			tt.validate(t, claims, err)
		})
	}
}

func TestTokenTiming(t *testing.T) {
	setupTestEnv(t)

	user := &GoogleUserInfo{
		ID:            "123",
		Email:         "test@test.com",
		VerifiedEmail: true,
	}

	token, err := GenerateToken(user)
	assert.NoError(t, err)

	claims, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	now := time.Now()

	// Test timestamps with wiggle room for execution time
	assert.True(t, claims.IssuedAt.Time.Before(now.Add(10*time.Second)))
	assert.True(t, claims.IssuedAt.Time.After(now.Add(-10*time.Second)))

	expectedExpiration := now.Add(24 * time.Hour)
	assert.True(t, claims.ExpiresAt.Time.Before(expectedExpiration.Add(10*time.Second)))
	assert.True(t, claims.ExpiresAt.Time.After(expectedExpiration.Add(-10*time.Second)))
}
