package auth

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/northwindman/testREST-autentification/internal/lib/format"
	"github.com/northwindman/testREST-autentification/internal/lib/logger/sl"
	"github.com/northwindman/testREST-autentification/internal/lib/random"
	"github.com/northwindman/testREST-autentification/internal/lib/tokens"
	"github.com/northwindman/testREST-autentification/internal/storage"
	"io"
	"log/slog"
	"net"
	"net/http"

	resp "github.com/northwindman/testREST-autentification/internal/lib/api/response"
)

type Request struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type Response struct {
	resp.Response
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserSaver interface {
	SaveUser(ip string, email string, passHash []byte, secret string, refreshToken []byte) (int64, error)
}

func New(log *slog.Logger, userSaver UserSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.auth.New"

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
			var validateErr validator.ValidationErrors
			if errors.As(err, &validateErr) {
				log.Error("invalid request", sl.Err(err))
				render.JSON(w, r, resp.ValidationError(validateErr))
			} else {
				log.Error("unexpected error", sl.Err(err))
				render.JSON(w, r, resp.Error("internal server error"))
			}
			return
		}

		secret, err := random.NewSecret(random.SecretLength)
		if err != nil {
			log.Error("failed to generate secret", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to generate secret"))
			return
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Error("failed to parse remote address", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to parse remote address"))
			return
		}

		token, err := tokens.GenTokens(ip, req.Email, secret, tokens.AccessTokenLength)
		if err != nil {
			log.Error("failed to generate token", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to generate token"))
			return
		}

		log.Info("generated token")

		tokenHash, err := format.HashString(token.RefreshToken)
		if err != nil {
			log.Error("failed to hash token", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to hash token"))
			return
		}

		token.RefreshToken = format.InBase64(token.RefreshToken)

		passHash, err := format.HashString(req.Password)
		if err != nil {
			log.Error("failed to hash password", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to hash password"))
			return
		}

		id, err := userSaver.SaveUser(ip, req.Email, passHash, secret, tokenHash)
		if errors.Is(err, storage.ErrAlreadyExist) {
			log.Warn("user already exists", sl.Err(err))
			render.JSON(w, r, resp.Error("user already exists"))
			return
		}
		if err != nil {
			log.Error("failed to save user", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to save user"))
			return
		}

		log.Info("user saved", slog.Int64("user", id))

		responseOK(w, r, token.AccessToken, token.RefreshToken)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, acToken string, rfToken string) {
	render.JSON(w, r, Response{
		Response:     resp.OK(),
		AccessToken:  acToken,
		RefreshToken: rfToken,
	})
}
