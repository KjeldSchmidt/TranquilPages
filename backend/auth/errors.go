package auth

import "fmt"

// OAuth flow errors
type (
	StateGenerationError   struct{ Err error }
	AuthURLGenerationError struct{ Err error }
	StateValidationError   struct{ Err error }
	TokenExchangeError     struct{ Err error }
	UserInfoError          struct{ Err error }
)

func (e *StateGenerationError) Error() string {
	return fmt.Sprintf("failed to generate state for OAuth flow: %v", e.Err)
}
func (e *AuthURLGenerationError) Error() string {
	return fmt.Sprintf("failed to generate auth URL: %v", e.Err)
}
func (e *StateValidationError) Error() string {
	return fmt.Sprintf("failed to validate state: %v", e.Err)
}
func (e *TokenExchangeError) Error() string {
	return fmt.Sprintf("failed to exchange token: %v", e.Err)
}
func (e *UserInfoError) Error() string { return fmt.Sprintf("failed to get user info: %v", e.Err) }
