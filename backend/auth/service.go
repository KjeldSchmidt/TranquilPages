package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"

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
	config    *oauth2.Config
	stateRepo *OAuthStateRepository
	stopChan  chan struct{}
	wg        sync.WaitGroup
}

func NewAuthService(config *oauth2.Config, stateRepo *OAuthStateRepository) *AuthService {
	service := &AuthService{
		config:    config,
		stateRepo: stateRepo,
		stopChan:  make(chan struct{}),
	}
	service.wg.Add(1)
	go service.startCleanupTask()
	return service
}

func (s *AuthService) Shutdown() {
	close(s.stopChan)
	s.wg.Wait()
}

func (s *AuthService) startCleanupTask() {
	defer s.wg.Done()
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopChan:
			return
		case <-ticker.C:
			if err := s.stateRepo.CleanupExpired(); err != nil {
				fmt.Printf("Failed to cleanup expired states: %v\n", err)
			}
		}
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
		return nil, fmt.Errorf("failed to validate state: %v", err)
	}
	if oauthState == nil {
		return nil, fmt.Errorf("invalid or expired state")
	}

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
