package auth

import "golang.org/x/oauth2"

// NewDropboxProvider defines details needed use Dropbox for OAuth2
// authorization.
//
// https://www.dropbox.com/developers/reference/oauth-guide
func NewDropboxProvider() *AuthProvider {
	return &AuthProvider{
		ID:      Dropbox,
		Key:     "dropbox",
		Enabled: true,
		Config: oauth2.Config{
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.dropbox.com/oauth2/authorize",
				TokenURL: "https://api.dropbox.com/oauth2/token",
			},
		},
	}
}
