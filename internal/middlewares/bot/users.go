package bot

import (
	"errors"

	"github.com/iskanye/mirea-queue/internal/models"
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

func (b *Bot) GetPermissions(handler telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		// Подразумевается что пользователь уже гарантировано существует
		user := c.Get("user").(models.User)
		if !user.QueueAccess {
			return c.Send("Вы не имете доступ к очереди")
		}

		return handler(c)
	}
}
