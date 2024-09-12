package refresh

import (
	"github.com/northwindman/testREST-autentification/internal/domain/models"
	resp "github.com/northwindman/testREST-autentification/internal/lib/api/response"
	"log/slog"
	"net/http"
)

type Response struct {
	AccessToken  string `json:"access_token" validate:"required"`
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type Request struct {
	resp.Response
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserProvider interface {
	GetUser(email string, secret string) (models.User, error)
	UpdateUser(email string, ip string, secret string, refreshToken []byte) (int64, error)
}

func New(log *slog.Logger, userProvider UserProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.user.refresh.New"

		log.With(
			slog.String("op", op),
		)

	}
}
