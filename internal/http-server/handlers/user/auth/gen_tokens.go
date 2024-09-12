package auth

import (
	"github.com/northwindman/testREST-autentification/internal/domain/models"
	"github.com/northwindman/testREST-autentification/internal/lib/tokens/access"
	"github.com/northwindman/testREST-autentification/internal/lib/tokens/jwt"
)

func genTokens(ip string, email string, secret string, accessTokenLength int) (models.Token, error) {
	rfToken, err := access.New(accessTokenLength)
	if err != nil {
		return models.Token{}, err
	}

	acToken, err := myjwt.New(ip, email, secret)
	if err != nil {
		return models.Token{}, err
	}

	return models.Token{
		AccessToken:  acToken,
		RefreshToken: rfToken,
	}, nil
}
