package tokens

import (
	"github.com/northwindman/testREST-autentification/internal/domain/models"
	"github.com/northwindman/testREST-autentification/internal/lib/tokens/jwt"
	"github.com/northwindman/testREST-autentification/internal/lib/tokens/refresh"
)

const (
	AccessTokenLength = 30
)

func GenTokens(ip string, email string, secret string, accessTokenLength int) (models.Token, error) {
	const op = "internal.lib.tokens.GenTokens"

	rfToken, err := refresh.New(accessTokenLength)
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
