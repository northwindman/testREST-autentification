package refresh

import (
	"crypto/rand"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// New create new token in format base64
func New(length int) (string, error) {
	const op = "lib.tokens.refresh.New"

	tokenBytes := make([]byte, length)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return string(tokenBytes), nil
}

// HashString hashes the input string
func HashString(incoming string) ([]byte, error) {
	const op = "lib.tokens.refresh.HashString"

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
