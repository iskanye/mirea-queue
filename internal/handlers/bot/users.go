package bot

import "gopkg.in/telebot.v4"

func (b *Bot) Start(c telebot.Context) error {
	return c.Send("hello")
}
