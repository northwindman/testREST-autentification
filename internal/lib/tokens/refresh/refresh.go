package refresh

import (
	"crypto/rand"
	"errors"
	"fmt"
)

var (
	ErrInvalidTokenLength = errors.New("invalid token length")
)

// New create new token string
func New(length int) (string, error) {
	const op = "lib.tokens.refresh.New"

	if length <= 0 {
		return "", fmt.Errorf(op+": %w", ErrInvalidTokenLength)
	}

	tokenBytes := make([]byte, length)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return string(tokenBytes), nil
}
