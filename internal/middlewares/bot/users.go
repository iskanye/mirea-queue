package bot

import (
	"errors"

	"github.com/iskanye/mirea-queue/internal/services"
	"gopkg.in/telebot.v4"
)

func (b *Bot) GetUser(handler telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		user, err := b.usersService.GetUser(b.ctx, c.Chat().ID)
		if errors.Is(err, services.ErrNotFound) {
			return c.Send("Вы не зарегистрированы в системе")
		}
		if err != nil {
			return nil
		}

		c.Set("user", user)
		return handler(c)
	}
}
