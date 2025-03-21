package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockTokenRepository implements TokenRepositoryInterface for testing
type MockTokenRepository struct {
	blacklistedTokens map[string]bool
	blacklistFunc     func(token string) error
}

func NewMockTokenRepository() *MockTokenRepository {
	return &MockTokenRepository{
		blacklistedTokens: make(map[string]bool),
	}
}

func (m *MockTokenRepository) Reset() {
	m.blacklistedTokens = make(map[string]bool)
	m.blacklistFunc = nil
}

func (m *MockTokenRepository) Blacklist(token string) error {
	if m.blacklistFunc != nil {
		return m.blacklistFunc(token)
	}
	m.blacklistedTokens[token] = true
	return nil
}

func (m *MockTokenRepository) IsBlacklisted(token string) (bool, error) {
	return m.blacklistedTokens[token], nil
}

// MockOAuthStateRepository implements OAuthStateRepositoryInterface for testing
type MockOAuthStateRepository struct {
	states     map[string]*OAuthState
	createFunc func(state *OAuthState) error
}

func NewMockOAuthStateRepository() *MockOAuthStateRepository {
	m := &MockOAuthStateRepository{
		states: make(map[string]*OAuthState),
	}
	return m
}

func (m *MockOAuthStateRepository) Reset() {
	m.states = make(map[string]*OAuthState)
	m.createFunc = nil
}

// Create stores a new OAuth state
func (m *MockOAuthStateRepository) Create(state *OAuthState) error {
	if m.createFunc != nil {
		return m.createFunc(state)
	}

	now := time.Now()
	state.CreatedAt = primitive.NewDateTimeFromTime(now)
	state.ExpiresAt = primitive.NewDateTimeFromTime(now.Add(15 * time.Minute))

	m.states[state.State] = state
	return nil
}

func (m *MockOAuthStateRepository) FindAndDelete(state string) (*OAuthState, error) {
	if s, exists := m.states[state]; exists {
		delete(m.states, state)
		return s, nil
	}
	return nil, nil
}
