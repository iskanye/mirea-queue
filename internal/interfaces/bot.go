package interfaces

import "gopkg.in/telebot.v4"

type BotHandlers interface {
	OnText(telebot.Context) error

	// Обработчики пользователей
	Start(telebot.Context) error
	Edit(telebot.Context) error

	// Обработчики очереди
	Push(telebot.Context) error
	Pop(telebot.Context) error
	LetAhead(telebot.Context) error
}

type BotMiddlewares interface {
	GetUser(telebot.HandlerFunc) telebot.HandlerFunc
	GetPermissions(telebot.HandlerFunc) telebot.HandlerFunc
}
