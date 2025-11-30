package interfaces

import "gopkg.in/telebot.v4"

type BotHandlers interface {
	OnText(c telebot.Context) error
	Start(c telebot.Context) error
}
