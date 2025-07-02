package oauth

// UserInfo represents basic user information returned by an OAuth provider.
type UserInfo struct {
	ID    string
	Email string
}

// Provider defines methods required for OAuth providers.
type Provider interface {
	// Name returns the provider name (e.g., "google" or "github").
	Name() string
	// LoginURL returns the authorization URL to redirect the user to.
	LoginURL(state string) string
	// Exchange converts an OAuth code into user information.
	Exchange(code string) (*UserInfo, error)
}
