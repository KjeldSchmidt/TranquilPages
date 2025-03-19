package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/oauth2"
)

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type AuthService struct {
	config *oauth2.Config
}

func NewAuthService(config *oauth2.Config) *AuthService {
	return &AuthService{
		config: config,
	}
}

// GetAuthURL generates the OAuth2 authorization URL
func (s *AuthService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state)
}

// HandleCallback processes the OAuth2 callback and returns user info
func (s *AuthService) HandleCallback(code string) (*GoogleUserInfo, error) {
	token, err := s.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %v", err)
	}

	client := s.config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %v", err)
	}
	defer resp.Body.Close()

	userInfo, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading user info: %v", err)
	}

	var user GoogleUserInfo
	if err := json.Unmarshal(userInfo, &user); err != nil {
		return nil, fmt.Errorf("failed parsing user info: %v", err)
	}

	return &user, nil
}
