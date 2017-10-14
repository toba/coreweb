package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	googleUserProfile = "https://www.googleapis.com/auth/userinfo.profile"
	googleUserEmail   = "https://www.googleapis.com/auth/userinfo.email"
)

// https://developers.google.com/identity/sign-in/web/
// https://developers.google.com/identity/sign-in/web/backend-auth
func NewGoogleProvider() *AuthProvider {
	return &AuthProvider{
		ID:      Google,
		Key:     "google",
		Enabled: true,
		Config: oauth2.Config{
			Scopes:   []string{googleUserProfile, googleUserEmail},
			Endpoint: google.Endpoint,
		},
	}
}
