package bot

import (
	"context"
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/interfaces"
)

type Bot struct {
	log *slog.Logger
	ctx context.Context

	queueService interfaces.QueueService
	usersService interfaces.UsersService
	adminService interfaces.AdminService
}

func New(
	log *slog.Logger,
	ctx context.Context,

	queueService interfaces.QueueService,
	usersService interfaces.UsersService,
	adminService interfaces.AdminService,
) *Bot {
	return &Bot{
		log: log,
		ctx: ctx,

		queueService: queueService,
		usersService: usersService,
		adminService: adminService,
	}
}
