package auth

// Repository interfaces
type TokenRepositoryInterface interface {
	Blacklist(token string) error
	IsBlacklisted(token string) (bool, error)
}

type OAuthStateRepositoryInterface interface {
	Create(state *OAuthState) error
	FindAndDelete(state string) (*OAuthState, error)
}
