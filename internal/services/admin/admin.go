package admin

import (
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/config"
)

type Admin struct {
	log   *slog.Logger
	token string
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *Admin {
	return &Admin{
		log:   log,
		token: cfg.AdminToken,
	}
}

func (a *Admin) ValidateToken(token string) bool {
	return token == a.token
}
