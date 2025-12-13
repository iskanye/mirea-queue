package bot

import (
	"errors"

	"github.com/iskanye/mirea-queue/internal/services"
	"gopkg.in/telebot.v4"
)

func (b *Bot) GetQueue(handler telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		queue, err := b.queueService.GetFromCache(b.ctx, c.Chat().ID)
		if errors.Is(err, services.ErrNotFound) {
			return c.Send("Предмет не найден, попробуйте указать его снова")
		} else if err != nil {
			return err
		}

		c.Set("queue", queue)
		return handler(c)
	}
}
