package myjwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

// NewToken creates a new JWT token for given user
func New(ip string, email string, secret string) (string, error) {
	const op = "lib.token.jwt.NewAccessToken"

	token := jwt.New(jwt.SigningMethodHS512)

	claims := token.Claims.(jwt.MapClaims)
	claims["ip"] = ip
	claims["email"] = email

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return tokenString, nil
}
