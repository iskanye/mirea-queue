package bot

import (
	"context"

	"github.com/iskanye/mirea-queue/internal/config"
	"github.com/iskanye/mirea-queue/internal/interfaces"
	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	b *tele.Bot

	startMenu *tele.ReplyMarkup
	editBtn   tele.Btn

	cancel context.CancelFunc
}

func New(
	cfg *config.Config,
) (*Bot, context.Context) {
	pref := tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: cfg.BotTimeout},
		OnError: func(err error, c tele.Context) {
			c.Send("Произошла неизвестная ошибка: " + err.Error())
		},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Инициализировать меню /start
	startMenu := &tele.ReplyMarkup{}
	edit := startMenu.Data("Изменить данные", "edit")
	startMenu.Inline(
		startMenu.Row(edit),
	)

	return &Bot{
		b:         b,
		startMenu: startMenu,
		cancel:    cancel,
	}, ctx
}

func (b *Bot) StartMenu() *tele.ReplyMarkup {
	return b.startMenu
}

func (b *Bot) Start() {
	b.b.Start()
}

func (b *Bot) Stop() {
	b.cancel()
	b.b.Stop()
}

func (b *Bot) Register(
	handlers interfaces.BotHandlers,
	middlewares interfaces.BotMiddlewares,
) {
	// Функции регистрации пользователя
	b.b.Handle("/start", handlers.Start)

	// Группа требующая авторизации
	authorized := b.b.Group()
	{
		authorized.Use(middlewares.GetUser)
		authorized.Handle("/edit", handlers.Edit)
		authorized.Handle("/push", handlers.Push)
		authorized.Handle("/swap", handlers.LetAhead)

		// Нужны права админа
		authorized.Handle("/pop", handlers.Pop, middlewares.GetPermissions)
	}

	// Обработчик любого текста
	b.b.Handle(tele.OnText, handlers.OnText)
}
