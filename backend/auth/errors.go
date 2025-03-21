package auth

import (
	"fmt"
)

// StateGenerationError represents an error that occurred while generating the OAuth state
type StateGenerationError struct {
	Err error
}

func (e *StateGenerationError) Error() string {
	return fmt.Sprintf("OAuth state generation error: %v", e.Err)
}

// AuthURLGenerationError represents an error that occurred while generating the OAuth URL
type AuthURLGenerationError struct {
	Err error
}

func (e *AuthURLGenerationError) Error() string {
	return fmt.Sprintf("OAuth auth URL generation error: %v", e.Err)
}

// StateValidationError represents an error that occurred while validating the OAuth state
type StateValidationError struct {
	Err error
}

func (e *StateValidationError) Error() string {
	return fmt.Sprintf("OAuth state validation error: %v", e.Err)
}

// TokenExchangeError represents an error that occurred during token exchange
type TokenExchangeError struct {
	Err error
}

func (e *TokenExchangeError) Error() string {
	return fmt.Sprintf("OAuth token exchange error: %v", e.Err)
}

// UserInfoError represents an error that occurred while fetching user info
type UserInfoError struct {
	Err error
}

func (e *UserInfoError) Error() string {
	return fmt.Sprintf("failed to get user info: %v", e.Err)
}

// TokenBlacklistError represents an error that occurred while managing token blacklist
type TokenBlacklistError struct {
	Err error
}

func (e *TokenBlacklistError) Error() string {
	return fmt.Sprintf("failed to manage token blacklist: %v", e.Err)
}
