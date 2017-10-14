package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/amazon"
)

// https://developer.amazon.com/public/apis/engage/login-with-amazon/docs/obtain_customer_profile.html#call_profile_endpoint
type amazonResponse struct {
	UserID       string `json:"user_id"`
	Name         string `json:"name"`
	PrimaryEmail string `json:"email"`
	PostalCode   string `json:"postal_code"`
}

// https://developer.amazon.com/public/apis/experience/cloud-drive/content/getting-started
const (
	amazonUserProfile = "profile"
)

// newAmazonProvider creates OAuth2 information to authorize using Amazon.
//
// https://developer.amazon.com/public/apis/engage/login-with-amazon/content/documentation.html
func NewAmazonProvider() *AuthProvider {
	return &AuthProvider{
		ID:              Amazon,
		Key:             "amazon",
		Enabled:         true,
		ProfileEndpoint: "https://api.amazon.com/user/profile",
		Config: oauth2.Config{
			Scopes:   []string{amazonUserProfile},
			Endpoint: amazon.Endpoint,
		},
	}
}
