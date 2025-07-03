package config

// OAuth provider endpoints and credentials
// Replace the placeholders with your actual client IDs and secrets.
var (
	GoogleClientID     = "GOOGLE_CLIENT_ID"
	GoogleClientSecret = "GOOGLE_CLIENT_SECRET"
	GoogleRedirectURL  = "http://localhost:8080/forum/api/oauth/google/callback"

	GitHubClientID     = "GITHUB_CLIENT_ID"
	GitHubClientSecret = "GITHUB_CLIENT_SECRET"
	GitHubRedirectURL  = "http://localhost:8080/forum/api/oauth/github/callback"
)
