package auth

import (
	"encoding/base64"
	"time"

	"github.com/toba/coreweb/token"
)

// State represents the OAuth2 state, encoded and passed with the standard
// authorization link and sent back to be decoded.
//
// https://tools.ietf.org/html/draft-bradley-oauth-jwt-encoded-state-05
type State struct {
	token.Token
	// whether login should also register the tenant
	Register bool
	// path to redirect to after login callback
	ModulePath    string
	AuthProviders []AuthProviderID
}

// NewState creates a token used for OAuth2 state.
func NewState(tenantID int64, modulePath string, authProviders ...AuthProviderID) *State {
	return &State{
		Token:         token.Base(tenantID, time.Duration(time.Hour)),
		ModulePath:    modulePath,
		AuthProviders: authProviders,
	}
}

// Encode serializes, signs and encodes a payload.
func (s *State) Encode() (string, error) {
	gob, err := token.ToGob(s)
	if err != nil {
		return "", err
	}
	hmac, err := token.Encode(gob)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(hmac), err
}

// DecodeState converts a string back to a state token with optional
// expiration and signature validation.
func DecodeState(state string, validate bool) (*State, error) {
	decoder, err := token.DecodeString(state, validate)
	if err != nil {
		return nil, err
	}
	s := &State{}
	err = decoder.Decode(s)

	if err != nil {
		return nil, err
	}
	if validate && s.IsExpired() {
		return nil, token.ErrTokenExpired
	}
	return s, nil
}
