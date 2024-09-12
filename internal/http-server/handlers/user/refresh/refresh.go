package refresh

import (
	"encoding/base64"
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/northwindman/testREST-autentification/internal/domain/models"
	resp "github.com/northwindman/testREST-autentification/internal/lib/api/response"
	"github.com/northwindman/testREST-autentification/internal/lib/logger/sl"
	"github.com/northwindman/testREST-autentification/internal/lib/notifications/email"
	"github.com/northwindman/testREST-autentification/internal/lib/random"
	"github.com/northwindman/testREST-autentification/internal/lib/tokens"
	myjwt "github.com/northwindman/testREST-autentification/internal/lib/tokens/jwt"
	"github.com/northwindman/testREST-autentification/internal/lib/tokens/refresh"
	"github.com/northwindman/testREST-autentification/internal/storage"
	"io"
	"log/slog"
	"net/http"
)

type Request struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type Response struct {
	resp.Response
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserProvider interface {
	GetUser(email string) (models.User, error)
	UpdateUser(email string, ip string, secret string, refreshToken []byte) (int64, error)
}

func New(log *slog.Logger, userProvider UserProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.refresh.New"

		log.With(
			slog.String("op", op),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", sl.Err(err))
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("failed to parse request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to parse request"))
			return
		}

		log.Info("request body decoded")

		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request")
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		claims, err := myjwt.GetClaims(req.AccessToken)
		if err != nil {
			log.Error("failed to get claims", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to get claims"))
			return
		}

		incomingEmail := claims["email"].(string)

		originalUser, err := userProvider.GetUser(incomingEmail)
		if err != nil {
			if errors.Is(err, storage.ErrNotFound) {
				log.Warn("user not found", sl.Err(err))
				render.JSON(w, r, resp.Error("user not found"))
				return
			}

			render.JSON(w, r, resp.Error("failed to get user"))
			return
		}

		decodedBytes, err := base64.StdEncoding.DecodeString(req.RefreshToken)
		if err != nil {
			log.Error("failed to decode refresh token", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode refresh token"))
			return
		}

		if ok := refresh.VerifyString(string(decodedBytes), originalUser.RefreshToken); !ok {
			log.Error("invalid refresh token")
			render.JSON(w, r, resp.Error("invalid credentials"))
			return
		}

		originalSecret := originalUser.Secret

		incomingUser, err := myjwt.ParseToken(req.AccessToken, originalSecret)
		if err != nil {
			log.Error("failed to parse token", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to parse token"))
			return
		}

		if originalUser.IP != incomingUser.IP {
			err = email.New(originalUser.Email, "Someone an another IP refresh your access token", "some body...")
			if err != nil {
				log.Error("failed to send email", sl.Err(err))
			}
		}

		newSecret, err := random.NewSecret(random.SecretLength)
		if err != nil {
			log.Error("failed to generate new secret", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		newTokens, err := tokens.GenTokens(incomingUser.IP, originalUser.Email, newSecret, tokens.AccessTokenLength)
		if err != nil {
			log.Error("failed to generate new tokens", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		tokenHash, err := refresh.HashString(newTokens.RefreshToken)
		if err != nil {
			log.Error("failed to hash token", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		newTokens.RefreshToken = base64.StdEncoding.EncodeToString([]byte(newTokens.RefreshToken))

		id, err := userProvider.UpdateUser(originalUser.Email, incomingUser.IP, newSecret, tokenHash)
		if err != nil {
			log.Error("failed to update user", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("user updated", slog.Int64("id", id))

		responseOK(w, r, newTokens.AccessToken, newTokens.RefreshToken)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, acToken string, rfToken string) {
	render.JSON(w, r, Response{
		Response:     resp.OK(),
		AccessToken:  acToken,
		RefreshToken: rfToken,
	})
}
