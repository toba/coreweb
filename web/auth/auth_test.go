package auth_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/toba/coreweb/web/auth"
)

var (
	key    = "TestKey"
	code   = "TestCode"
	state  = "TestState"
	secret = "TestSecret"
)

func TestAuthInit(t *testing.T) {
	p := auth.GetProvider(auth.Dropbox)
	assert.NotNil(t, p)

	p.Initialize(key, secret)
	assert.Equal(t, key, p.ClientID)
	assert.Equal(t, secret, p.ClientSecret)
}

func TestAuthLink(t *testing.T) {
	p := auth.GetProvider(auth.Dropbox)
	assert.NotNil(t, p)

	p.Initialize(key, secret)

	redirect := "http://base-url.com"
	url := p.Endpoint.AuthURL +
		"?client_id=" + key +
		"&redirect_uri=" + url.QueryEscape(redirect+"/auth/"+p.Key) +
		"&response_type=code&state=" + state

	authLink := p.Link(redirect, state)
	assert.NotNil(t, authLink)
	assert.Equal(t, auth.Dropbox, authLink.ID)
	assert.Equal(t, p.Key, authLink.Key)
	assert.Equal(t, url, authLink.URL)
}

func TestCallbackHandler(t *testing.T) {
	p := auth.GetProvider(auth.Dropbox)
	assert.NotNil(t, p)

	qs := url.Values{}
	qs.Set("code", code)
	qs.Add("state", state)

	h := http.HandlerFunc(p.HandleCallback)

	server := httptest.NewServer(h)
	defer server.Close()

	res, err := http.Get(server.URL + "?" + qs.Encode())

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}
