// Package token is a simplified implementation of Javascript Web Tokens
// (JWT).
//
// The JWT header segment has been removed since the GWT hash algorithm
// and type do not change.
//
// HMAC is keyed-Hash Message Authentication Code. Different from a digital
// signature, it uses the same key to sign and verify a payload, meaning the
// key must be shared or the payload can only be validated by the same
// process that generates it.
//
// A token cannot be modified without it failing verification. Data it
// contains are not guaranteed to remain private but are guaranteed
// to remain unmodified as long as the key is secure.
//
// See https://blake2.net/
//
package token

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"math"
	"time"

	"golang.org/x/crypto/blake2b"
)

type (
	Token struct {
		TenantID int64
		Expires  time.Time
	}
)

var (
	key                 = []byte("SnrKNv(5zecfZLYc0nGhG9OPnM$dtz4H7hm@SsxIlhJCp4LOa*LWuLo")
	ErrTokenExpired     = errors.New("token expired")
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidSignature = errors.New("signature verification failed")
	ErrTooLargeForToken = errors.New("payload too large for token")
)

// IsExpired indicates whether token payload is expired.
func (t *Token) IsExpired() bool {
	return t.Expires.Before(time.Now())
}

// // Encode serializes and signs a payload.
// func (t *Token) Encode() ([]byte, error) {
// 	gob, err := ToGob(t)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return encode(gob)
// }

// func (t *Token) EncodeAsString() (string, error) {
// 	gob, err := t.Encode()
// 	if err != nil {
// 		return "", err
// 	}
// 	return base64.URLEncoding.EncodeToString(gob), err
// }

func Encode(gob []byte) ([]byte, error) {
	if len(gob) > math.MaxUint16 {
		return nil, ErrTooLargeForToken
	}

	signature, err := sign(gob)
	if err != nil {
		return nil, err
	}

	size := make([]byte, 2)
	binary.LittleEndian.PutUint16(size, uint16(len(gob)+2))

	token := size[:]
	token = append(token, gob...)

	return append(token, signature...), nil
}

// ToGob encodes the payload using gob.
func ToGob(t interface{}) ([]byte, error) {
	var buf bytes.Buffer
	en := gob.NewEncoder(&buf)
	if err := en.Encode(t); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func sign(encoded []byte) ([]byte, error) {
	hash, err := blake2b.New256(key)
	if err != nil {
		return nil, err
	}
	if _, err := hash.Write(encoded); err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

// verify ensures that a signature is correct for the encoded bytes.
func verify(encoded, signature []byte) bool {
	testSign, err := sign(encoded)
	if err != nil {
		return false
	}
	return bytes.Equal(testSign, signature)
}

// Decode converts a hashed token to a token structure. The hash is
// formatted as <header>.<payload>.<signature>
func Decode(token []byte, validate bool) (*Token, error) {
	decoder, err := decode(token, validate)
	if err != nil {
		return nil, err
	}
	t := &Token{}
	err = decoder.Decode(t)

	if err != nil {
		return nil, err
	}

	if validate && t.IsExpired() {
		return nil, ErrTokenExpired
	}
	return t, nil
}

// decode converts a hashed token to a token structure. The hash is
// formatted as <header>.<payload>.<signature>
func decode(token []byte, validate bool) (*gob.Decoder, error) {
	if token == nil || len(token) < 5 {
		return nil, ErrInvalidToken
	}
	var buf bytes.Buffer
	size := binary.LittleEndian.Uint16(token[0:2])
	encoded := token[2:size]
	signature := token[size:]

	decoder := gob.NewDecoder(&buf)
	_, err := buf.Write(encoded)

	if err != nil {
		return nil, err
	}

	if validate && !verify(encoded, signature) {
		return nil, ErrInvalidSignature
	}
	return decoder, err
}

func DecodeString(token string, validate bool) (*gob.Decoder, error) {
	hmac, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	return decode(hmac, validate)
}

// Base creates the basic token struct.
func Base(tenantID int64, duration time.Duration) Token {
	return Token{
		TenantID: tenantID,
		Expires:  time.Now().Add(duration),
	}
}
