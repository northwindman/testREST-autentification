package myjwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/northwindman/testREST-autentification/internal/domain/models"
)

var (
	ErrEmptyClaims   = errors.New("empty claims")
	ErrInvalidClaims = errors.New("invalid claims")
)

// NewToken creates a new JWT token for given user
func New(ip string, email string, secret string) (string, error) {
	const op = "lib.token.jwt.NewAccessToken"

	if secret == "" {
		return "", fmt.Errorf("empty secret")
	}

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
		ip, ipOk := claims["ip"].(string)
		if !ipOk || ip == "" {
			return nil, fmt.Errorf("%s: %w", op, ErrEmptyClaims)
		}
		email, emailOk := claims["email"].(string)
		if !emailOk || email == "" {
			return nil, fmt.Errorf("%s: %w", op, ErrEmptyClaims)
		}

		return claims, nil
	}

	return nil, fmt.Errorf("%s: %w", op, ErrInvalidClaims)
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
		if ip, ipOk := claims["ip"].(string); ipOk {
			user.IP = ip
		} else {
			return models.User{}, fmt.Errorf("%s: invalid or missing 'ip' claim", op)
		}

		if email, emailOk := claims["email"].(string); emailOk {
			user.Email = email
		} else {
			return models.User{}, fmt.Errorf("%s: invalid or missing 'email' claim", op)
		}

	} else {
		return models.User{}, fmt.Errorf("%s: %w", op, ErrInvalidClaims)
	}

	return user, nil
}
