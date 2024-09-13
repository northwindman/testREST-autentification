package format

import (
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmptyValue = errors.New("empty value")
)

// HashString hashes the input string
func HashString(incoming string) ([]byte, error) {
	const op = "lib.tokens.refresh.HashString"

	if len(incoming) == 0 {
		return []byte{}, ErrEmptyValue
	}

	hashedString, err := bcrypt.GenerateFromPassword([]byte(incoming), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return hashedString, nil
}

// VerifyString compares the received value with the hash
func VerifyString(received string, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(received))
	return err == nil
}

// InBase64 return string in format base64
func InBase64(received string) string {
	return base64.StdEncoding.EncodeToString([]byte(received))
}

// FromBase64 return default string
func FromBase64(received string) (string, error) {
	const op = "lib.tokens.refresh.FromBase64"

	if len(received) == 0 {
		return "", ErrEmptyValue
	}

	b, err := base64.StdEncoding.DecodeString(received)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return string(b), nil
}
