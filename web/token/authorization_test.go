package token_test

import (
	"testing"

	"github.com/toba/goweb/web/token"

	"github.com/stretchr/testify/assert"
)

func TestAuthorizationEncoding(t *testing.T) {
	key := int64(123)
	p := token.ForAuth(key, 1, 2, 3, 4)

	enc, err := p.Encode()
	assert.NoError(t, err)
	assert.NotNil(t, enc)

	dec, err := token.DecodeAuthorization(enc, true)
	assert.NoError(t, err)
	assert.NotNil(t, dec)
	assert.Equal(t, key, dec.TenantID)
	assert.Contains(t, dec.Permissions, uint16(3))
}
