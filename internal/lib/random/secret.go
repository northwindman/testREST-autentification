package random

import (
	"crypto/rand"
	"encoding/base64"
)

const (
	SecretLength = 32
)

// NewSecret returns a string for signing the token
func NewSecret(length int) (string, error) {
	byteArr := make([]byte, length)
	_, err := rand.Read(byteArr)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(byteArr)[:length], nil
}
