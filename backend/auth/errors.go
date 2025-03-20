package auth

import "fmt"

// OAuth flow errors
type (
	StateGenerationError   struct{ Err error }
	AuthURLGenerationError struct{ Err error }
	StateValidationError   struct{ Err error }
	TokenExchangeError     struct{ Err error }
	UserInfoError          struct{ Err error }
	TokenBlacklistError    struct{ Err error }
)

func (e *StateGenerationError) Error() string {
	return fmt.Sprintf("OAuth state generation error: %v", e.Err)
}
func (e *AuthURLGenerationError) Error() string {
	return fmt.Sprintf("OAuth auth URL generation error: %v", e.Err)
}
func (e *StateValidationError) Error() string {
	return fmt.Sprintf("OAuth state validation error: %v", e.Err)
}
func (e *TokenExchangeError) Error() string {
	return fmt.Sprintf("OAuth token exchange error: %v", e.Err)
}
func (e *UserInfoError) Error() string       { return fmt.Sprintf("user info error: %v", e.Err) }
func (e *TokenBlacklistError) Error() string { return fmt.Sprintf("token blacklist error: %v", e.Err) }
