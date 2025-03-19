package auth

import "fmt"

type StateGenerationError struct {
	Err error
}

func (e *StateGenerationError) Error() string {
	return fmt.Sprintf("failed to generate state for OAuth flow: %v", e.Err)
}

type AuthURLGenerationError struct {
	Err error
}

func (e *AuthURLGenerationError) Error() string {
	return fmt.Sprintf("failed to generate auth URL: %v", e.Err)
}
