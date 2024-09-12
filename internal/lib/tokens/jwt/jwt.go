package myjwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/northwindman/testREST-autentification/internal/domain/models"
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

// GetClaims returns jwt.MapClaims for check fields
func GetClaims(tokenString string) (jwt.MapClaims, error) {
	const op = "lib.token.jwt.GetClaims"

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("%s: map.Claims is empty", op)
}

// ParseToken check if token is valid and original
func ParseToken(tokenString string, secret string) (models.User, error) {
	const op = "lib.token.jwt.Parse"

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if token.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	var user models.User

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user.Email = claims["email"].(string)
		user.IP = claims["ip"].(string)
	} else {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
