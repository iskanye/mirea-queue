package interfaces

import "gopkg.in/telebot.v4"

type BotHandlers interface {
	Start(c telebot.Context) error
}
