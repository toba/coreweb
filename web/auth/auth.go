// Package auth manages OAuth2 authentication. The flow applied is code-based
// instead of implicit meaning the steps are:
//
// - Random ID generated in browser and saved to cookie `authState`. Cookies are
//   used becuase they're automatically passed between browser and server in the
//   header.
// - Current browser path saved to cookie `authPath`.
// - The random ID is also sent with a service call to list authentication
//   providers so that the providers' generated redirect URLs include it as
//   their standard OAuth state value.
// - User clicks on one of the resulting provider links and authenticates.
// - Provider redirects browser to the redirect URL which includes the state
//   value and an authentication code.
// - Browser sends cookies in header when it loads the redirect URL allowing the
//   server to compare the original random cookie ID to the state value passed
//   back from the provider.
// - If state values match then server exchanges the code for token and writes
//	  token to cookie `<AuthProvider.Key>AccessToken`.
// - Server redirects browser to the view name saved in `authPath`.
// - Browser may then read and use the new access token.
//
// OAuth is designed for authorization, not authentication, which is why user
// identity isn't exposed to the requestor, only a temporary token. OpenID
// Connect is one approach to add authentication to OAuth by including a hash
// with user identity (such as JWT) in the token response.
//
// Most OAuth providers suggest a more ad hoc approach to authentication which
// is to use the received token to immediately make an API call to retrieve user
// identity. This approach is criticized as insecure or not true authentication
// since the requestor may retrieve user identity even after the user has
// stopped interacting with the application.
//
// https://oauth.net/articles/authentication/
package auth

import (
	"context"
	"log"
	"net/http"

	"golang.org/x/oauth2"
)

type (
	AuthProviderID int

	// AuthProvider describes how to authenticate against an OAuth2 provider.
	AuthProvider struct {
		oauth2.Config
		ID              AuthProviderID
		ProfileEndpoint string
		Enabled         bool
		// Key is a simple word to identify the provider's logo, description,
		// callback URL and the environment variable containing the OAuth key and
		// secret.
		Key string
	}

	// AuthLink encapsulates the per-user URL and values used by the React client
	// to initiate OAuth.
	AuthLink struct {
		ID  AuthProviderID `json:"id"`
		Key string         `json:"key"`
		URL string         `json:"url"`
	}
)

// Hardcode specific numbers instead of using iota to ensure they don't change.
const (
	None      AuthProviderID = 0
	LDAP      AuthProviderID = 10
	Amazon    AuthProviderID = 20
	Dropbox   AuthProviderID = 30
	Facebook  AuthProviderID = 40
	GitHub    AuthProviderID = 50
	Google    AuthProviderID = 60
	Microsoft AuthProviderID = 70
	PayPal    AuthProviderID = 80
	Slack     AuthProviderID = 90
	Twitter   AuthProviderID = 100
	Yahoo     AuthProviderID = 110
)

// list of active authentication providers.
var providers []*AuthProvider

// Initialize defines the active authentication providers and login callbacks.
func Initialize(p ...*AuthProvider) error {
	providers = p
	return nil
}

// GetProvider returns the Provider instance with the given ID.
func GetProvider(id AuthProviderID) *AuthProvider {
	for _, p := range providers {
		if p.ID == id {
			return p
		}
	}
	return nil
}

// Providers returns all enabled authentication providers.
func Providers() []*AuthProvider {
	enabled := []*AuthProvider{}
	for _, p := range providers {
		if p.Enabled {
			enabled = append(enabled, p)
		}
	}
	return enabled
}

// GetProviderForKey returns the Provider instance with the given Key.
func GetProviderForKey(key string) *AuthProvider {
	for _, p := range providers {
		if p.Key == key {
			return p
		}
	}
	return nil
}

// Initialize updates provider with client secrets and constructs the OAuth2
// endpoint.
func (auth *AuthProvider) Initialize(clientID, clientSecret string) {
	auth.ClientID = clientID
	auth.ClientSecret = clientSecret
}

// Link builds a simplified struct for browser client use.
func (auth *AuthProvider) Link(baseURL, state string) *AuthLink {
	auth.RedirectURL = baseURL + "/auth/" + auth.Key
	return &AuthLink{
		ID:  auth.ID,
		Key: auth.Key,
		URL: auth.AuthCodeURL(state),
	}
}

// // SetAccessToken sets access token to avoid calling Auth method.
// func (auth *AuthProvider) SetAccessToken(accessToken string) {
// 	auth.token = &oauth2.Token{AccessToken: accessToken}
// }

// // AccessToken returns the OAuth access token.
// func (auth *AuthProvider) AccessToken() string {
// 	return auth.token.AccessToken
// }

// func (auth *AuthProvider) client() *http.Client {
// 	return auth.Client(oauth2.NoContext, auth.token)
// }

// HandleCallback handles the provider's call to the response URI after user
// authenticates. It should convert the response code to a bearer token. See
// example at
//
// https://jacobmartins.com/2016/02/29/getting-started-with-oauth2-in-go/
func (auth *AuthProvider) HandleCallback(w http.ResponseWriter, r *http.Request) {
	var code string
	var state string
	var error string

	if r.Method == http.MethodGet {
		code = r.URL.Query().Get("code")
		state = r.URL.Query().Get("state")
		error = r.URL.Query().Get("error")
	} else if r.Method == http.MethodPost {
		code = r.FormValue("code")
		state = r.FormValue("state")
		error = r.FormValue("error")
	} else {
		http.Error(w, "", http.StatusMethodNotAllowed)
	}

	if error != "" {
		http.Error(w, error, http.StatusExpectationFailed)
		return
	}

	if code == "" || state == "" {
		http.Error(w, "Invalid authentication code or state", http.StatusUnauthorized)
		return
	}

	stateToken, err := DecodeState(state, true)
	if err != nil {
		http.Error(w, "Invalid authentication state", http.StatusUnauthorized)
		return
	}

	token, err := auth.Exchange(context.Background(), code)
	if err != nil {
		// TODO: log error and show friendly page instead of outputing error
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if stateToken.ModulePath == "" {
		http.Error(w, "No redirect path found", http.StatusNotFound)
		return
	}

	//err = tenantLogin(stateToken)

	log.Print("Access Token: " + token.AccessToken)

	http.SetCookie(w, &http.Cookie{Name: auth.Key + "AccessToken", Value: token.AccessToken, Path: "/"})
	http.Redirect(w, r, stateToken.ModulePath, http.StatusTemporaryRedirect)
}
