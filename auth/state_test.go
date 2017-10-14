package auth_test

import (
	"testing"

	"github.com/toba/coreweb/auth"

	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	key := int64(123)
	p := auth.NewState(key, "ModulePath", auth.Amazon, auth.GitHub)

	enc, err := p.Encode()
	assert.NoError(t, err)
	assert.NotNil(t, enc)

	dec, err := auth.DecodeState(enc, true)
	assert.NoError(t, err)
	assert.NotNil(t, dec)
	assert.Equal(t, key, dec.TenantID)
	assert.Contains(t, dec.AuthProviders, auth.Amazon)
}
