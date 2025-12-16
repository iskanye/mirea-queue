package bot

import (
	"log/slog"

	"gopkg.in/telebot.v4"
)

func (b *Bot) Logger(handler telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		log := b.log.With(
			slog.Int64("id", c.Chat().ID),
		)

		log.Info("COMMAND STARTED")
		defer log.Info("COMMANDED ENDED")

		return handler(c)
	}
}
