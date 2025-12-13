package bot

import (
	"context"
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/interfaces"
	"gopkg.in/telebot.v4"
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

func (b *Bot) CallbackRespond(h telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if c.Callback() != nil {
			defer c.Respond()
		}
		return h(c)
	}
}
