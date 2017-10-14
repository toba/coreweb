package auth

import "golang.org/x/oauth2"

// https://developer.yahoo.com/auth/
func NewYahooProvider() *AuthProvider {
	return &AuthProvider{
		ID:      Yahoo,
		Key:     "yahoo",
		Enabled: true,
		Config: oauth2.Config{
			Scopes: []string{googleUserProfile, googleUserEmail},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.facebook.com/v2.10/dialog/oauth",
				TokenURL: "https://api.dropbox.com/1/oauth2/token",
			},
		},
	}
}
