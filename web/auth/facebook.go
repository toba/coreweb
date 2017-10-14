package auth

import (
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

type facebookDebugResponse struct {
	AppID       int64     `json:"app_id"`
	UserID      int       `json:"user_id"`
	Application string    `json:"application"`
	ExpiresAt   time.Time `json:"expires_at"`
	IsValid     bool      `json:"is_valid"`
	IssuedAt    time.Time `json:"issued_at"`
}

// https://developers.facebook.com/docs/facebook-login/permissions/
const (
	fbPublicProfile = "public_profile"
	fbEMail         = "email"
)

// https://developers.facebook.com/docs/facebook-login/manually-build-a-login-flow
func NewFacebookProvider() *AuthProvider {
	return &AuthProvider{
		ID:      Facebook,
		Key:     "facebook",
		Enabled: true,
		Config: oauth2.Config{
			Endpoint: facebook.Endpoint,
		},
	}
}
