package oauth

import "context"

// UserInfo holds normalized user data returned from any OAuth provider.
type UserInfo struct {
	ProviderID string // Provider's unique user ID
	Email      string // May be empty (e.g. QQ)
	Name       string // Display name
	Avatar     string // Avatar URL
}

// Provider is the interface every OAuth provider must implement.
// To add a new provider: create a new file, implement this interface,
// and call Register() in an init() function.
type Provider interface {
	// Name returns the unique key, e.g. "github", "google", "qq".
	Name() string
	// AuthURL returns the authorization URL to redirect the user to.
	AuthURL(state string) string
	// Exchange trades an authorization code for normalized user info.
	Exchange(ctx context.Context, code string) (*UserInfo, error)
}
