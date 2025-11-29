package bot

import "gopkg.in/telebot.v4"

func (b *Bot) registerHandlers() {
	b.b.Handle("/start", func(c telebot.Context) error {
		name := c.Chat().FirstName

		return c.Send(name)
	})
}
