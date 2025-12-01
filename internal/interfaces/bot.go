package interfaces

import "gopkg.in/telebot.v4"

type BotHandlers interface {
	OnText(telebot.Context) error
	Start(telebot.Context) error
	Edit(telebot.Context) error
}
