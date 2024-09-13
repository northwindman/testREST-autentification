package myjwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test function New from this package

func TestNewAccessToken_Success(t *testing.T) {
	ip := "127.0.0.1"
	email := "test@example.com"
	secret := "mysecretkey"

	tokenString, err := New(ip, email, secret)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Check that token can parsed with same secret
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	// Check data in the token
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		assert.Equal(t, ip, claims["ip"])
		assert.Equal(t, email, claims["email"])
	} else {
		t.Fatal("claims not valid")
	}
}

func TestNewAccessToken_SigningError(t *testing.T) {
	ip := "127.0.0.1"
	email := "test@example.com"
	secret := "" // empty secret

	tokenString, err := New(ip, email, secret)

	assert.Error(t, err)
	assert.Equal(t, "", tokenString)
}

func TestNewAccessToken_InvalidSecret(t *testing.T) {
	ip := "127.0.0.1"
	email := "test@example.com"
	secret := "mysecretkey"
	invalidSecret := "wrongsecret"

	tokenString, err := New(ip, email, secret)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	// Parsing with invalid secret
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(invalidSecret), nil
	})

	assert.Error(t, err)
	assert.False(t, parsedToken.Valid)
}

func TestNewAccessToken_EmptyEmail(t *testing.T) {
	ip := "127.0.0.1"
	email := "" // Empty email
	secret := "mysecretkey"

	tokenString, err := New(ip, email, secret)

	assert.NoError(t, err)
	assert.NotEmpty(t, tokenString)

	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	// Check tokens data
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		assert.Equal(t, ip, claims["ip"])
		assert.Equal(t, email, claims["email"])
	} else {
		t.Fatal("claims not valid")
	}
}

// Test function GetClaims from this package

func TestGetClaims_Success(t *testing.T) {
	secret := "mysecretkey"
	tokenString := createValidToken(t, secret)

	claims, err := GetClaims(tokenString)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "127.0.0.1", claims["ip"])
	assert.Equal(t, "test@example.com", claims["email"])
}

func TestGetClaims_InvalidToken(t *testing.T) {
	invalidTokenString := "invalid.token.string"

	claims, err := GetClaims(invalidTokenString)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestGetClaims_EmptyClaims(t *testing.T) {
	emptyTokenString := createTokenWithEmptyClaims(t)

	claims, err := GetClaims(emptyTokenString)

	assert.Error(t, err)
	assert.Nil(t, claims)
}

func createValidToken(t *testing.T, secret string) string {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["ip"] = "127.0.0.1"
	claims["email"] = "test@example.com"

	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	return tokenString
}

func createTokenWithEmptyClaims(t *testing.T) string {
	token := jwt.New(jwt.SigningMethodHS512)
	token.Claims = jwt.MapClaims{}

	tokenString, err := token.SignedString([]byte("mysecretkey"))
	assert.NoError(t, err)

	return tokenString
}

// Test function ParseToken from this package

func TestParseToken_Success(t *testing.T) {
	secret := "mysecret"
	userEmail := "test@example.com"
	userIP := "192.168.1.1"

	tokenString, err := createTestToken(userEmail, userIP, secret)
	assert.NoError(t, err)

	parsedUser, err := ParseToken(tokenString, secret)
	assert.NoError(t, err)

	assert.Equal(t, userEmail, parsedUser.Email)
	assert.Equal(t, userIP, parsedUser.IP)
}

func TestParseToken_InvalidAlgorithm(t *testing.T) {
	secret := "mysecret"
	userEmail := "test@example.com"
	userIP := "192.168.1.1"

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = userEmail
	claims["ip"] = userIP

	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	_, err = ParseToken(tokenString, secret)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected signing method")
}

func TestParseToken_InvalidToken(t *testing.T) {
	secret := "mysecret"

	_, err := ParseToken("invalidTokenString", secret)
	assert.Error(t, err)
}

func TestParseToken_EmptyClaims(t *testing.T) {
	secret := "mysecret"

	token := jwt.New(jwt.SigningMethodHS512)
	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	_, err = ParseToken(tokenString, secret)
	assert.Error(t, err)
}

func TestParseToken_InvalidEmailIPType(t *testing.T) {
	secret := "mysecret"

	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = 123
	claims["ip"] = true

	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	_, err = ParseToken(tokenString, secret)
	assert.Error(t, err)
}

func createTestToken(email, ip, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["ip"] = ip

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
