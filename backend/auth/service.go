package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
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
	config      *oauth2.Config
	stateRepo   OAuthStateRepositoryInterface
	tokenRepo   TokenRepositoryInterface
	userInfoURL string
}

func NewAuthService(config *oauth2.Config, stateRepo OAuthStateRepositoryInterface, tokenRepo TokenRepositoryInterface) *AuthService {
	return &AuthService{
		config:      config,
		stateRepo:   stateRepo,
		tokenRepo:   tokenRepo,
		userInfoURL: "https://www.googleapis.com/oauth2/v2/userinfo",
	}
}

// GenerateRandomState generates a random state string for OAuth flow
func GenerateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", &StateGenerationError{Err: err}
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetAuthURL generates the OAuth2 authorization URL
func (s *AuthService) GetAuthURL() (string, error) {
	state, err := GenerateRandomState()
	if err != nil {
		return "", &AuthURLGenerationError{Err: fmt.Errorf("state parameter is required")}
	}

	oauthState := &OAuthState{
		State: state,
	}

	if err := s.stateRepo.Create(oauthState); err != nil {
		return "", &AuthURLGenerationError{Err: fmt.Errorf("failed to store state: %v", err)}
	}

	return s.config.AuthCodeURL(state), nil
}

// HandleCallback processes the OAuth2 callback and returns user info
func (s *AuthService) HandleCallback(code, state string) (*GoogleUserInfo, error) {
	oauthState, err := s.stateRepo.FindAndDelete(state)
	if err != nil {
		return nil, &StateValidationError{Err: err}
	}
	if oauthState == nil {
		return nil, &StateValidationError{Err: fmt.Errorf("invalid or expired state")}
	}

	token, err := s.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, &TokenExchangeError{Err: err}
	}

	client := s.config.Client(context.Background(), token)
	resp, err := client.Get(s.userInfoURL)
	if err != nil {
		return nil, &UserInfoError{Err: fmt.Errorf("failed getting user info: %v", err)}
	}
	defer resp.Body.Close()

	userInfo, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &UserInfoError{Err: fmt.Errorf("failed reading user info: %v", err)}
	}

	var user GoogleUserInfo
	if err := json.Unmarshal(userInfo, &user); err != nil {
		return nil, &UserInfoError{Err: fmt.Errorf("failed parsing user info: %v", err)}
	}

	return &user, nil
}

// Logout disables the given token, ensuring that it cannot be used for further authentication
func (s *AuthService) Logout(token string) error {
	_, err := ValidateToken(token)
	if err != nil {
		return err
	}

	// Blacklist the token until its expiration
	if err := s.tokenRepo.Blacklist(token); err != nil {
		return fmt.Errorf("failed to blacklist token: %v", err)
	}

	return nil
}

// ValidateAuthenticationToken checks if a token is valid and active, extracting Claims for further use if so.
func (s *AuthService) ValidateAuthenticationToken(tokenString string) (*Claims, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	// Check if token is blacklisted
	isBlacklisted, err := s.tokenRepo.IsBlacklisted(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to check token blacklist: %v", err)
	}
	if isBlacklisted {
		return nil, &TokenBlacklistError{Err: fmt.Errorf("token has been revoked")}
	}

	return claims, nil
}
