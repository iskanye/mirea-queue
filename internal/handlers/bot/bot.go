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

	channels map[int64]chan *telebot.Message

	queueService interfaces.QueueService
	usersService interfaces.UsersService
}

func New(
	log *slog.Logger,
	ctx context.Context,

	queueService interfaces.QueueService,
	usersService interfaces.UsersService,
) *Bot {
	return &Bot{
		log: log,
		ctx: ctx,

		channels: make(map[int64]chan *telebot.Message),

		queueService: queueService,
		usersService: usersService,
	}
}

// Обрабатывает текстовый ввод пользователя
func (b *Bot) OnText(c telebot.Context) error {
	if ch, ok := b.channels[c.Chat().ID]; ok {
		ch <- c.Message()
		return nil
	}

	b.log.Warn("Unhandled input",
		slog.String("text", c.Text()),
	)
	return nil
}

// Оборачивает команду в диалог, ожидая ввода от пользователя.
// Ввод пользователя передается через канал ch
func (b *Bot) Dialogue(
	c telebot.Context,
	fun func(ch <-chan *telebot.Message, c telebot.Context) error,
) error {
	chatID := c.Chat().ID
	ch := make(chan *telebot.Message, 1)
	b.channels[chatID] = ch

	defer func(cmd string) {
		close(ch)
		delete(b.channels, chatID)

		b.log.Info("Dialogue chain closed",
			slog.String("cmd", cmd),
		)
	}(c.Text())

	b.log.Info("Dialogue chain started",
		slog.String("cmd", c.Text()),
	)
	return fun(ch, c)
}
