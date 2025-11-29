package bot

import (
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/interfaces"
)

type Bot struct {
	log *slog.Logger

	queueService interfaces.QueueService
	usersService interfaces.UsersService
}

func New(
	log *slog.Logger,

	queueService interfaces.QueueService,
	usersService interfaces.UsersService,
) *Bot {
	return &Bot{
		log: log,

		queueService: queueService,
		usersService: usersService,
	}
}
