package token

import "time"

type AuthToken struct {
	Token
	Permissions []uint16
}

// ForAuthorization creates an authorization token.
func ForAuth(tenantID int64, permissions ...uint16) *AuthToken {
	return &AuthToken{
		Token:       Base(tenantID, time.Duration(time.Hour*24)),
		Permissions: permissions,
	}
}

// Encode serializes and signs a payload.
func (t *AuthToken) Encode() ([]byte, error) {
	gob, err := ToGob(t)
	if err != nil {
		return nil, err
	}
	return Encode(gob)
}

func DecodeAuthorization(token []byte, validate bool) (*AuthToken, error) {
	decoder, err := decode(token, validate)
	if err != nil {
		return nil, err
	}
	t := &AuthToken{}
	err = decoder.Decode(t)

	if err != nil {
		return nil, err
	}
	if validate && t.IsExpired() {
		return nil, ErrTokenExpired
	}
	return t, nil
}
