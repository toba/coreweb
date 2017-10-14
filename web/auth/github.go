package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

// https://developer.github.com/v3/guides/basics-of-authentication/
func NewGitHubProvider() *AuthProvider {
	return &AuthProvider{
		ID:      GitHub,
		Key:     "github",
		Enabled: true,
		Config: oauth2.Config{
			Scopes:   []string{googleUserProfile, googleUserEmail},
			Endpoint: github.Endpoint,
		},
	}
}
