package jwt

import "github.com/northwindman/testREST-autentification/internal/domain/models"

// NewToken creates a new JWT token for given user
func NewAccessToken(user models.User, app models.App) (string, error) {}
