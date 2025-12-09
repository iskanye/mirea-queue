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
	editBtn   *tele.Btn
	chooseBtn *tele.Btn

	subjectMenu *tele.ReplyMarkup
	refreshBtn  *tele.Btn
	pushBtn     *tele.Btn
	popBtn      *tele.Btn
	letAheadBtn *tele.Btn

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
	choose := startMenu.Data("Выбрать предмет", "choose")
	startMenu.Inline(
		startMenu.Row(edit),
		startMenu.Row(choose),
	)

	// Меню предмета
	subjectMenu := &tele.ReplyMarkup{}
	refresh := subjectMenu.Data("Обновить", "update")
	push := subjectMenu.Data("Записаться в очередь", "push")
	pop := subjectMenu.Data("Позвать на сдачу", "pop")
	letAhead := subjectMenu.Data("Пропустить в очереди", "let-ahead")
	subjectMenu.Inline(
		subjectMenu.Row(refresh),
		subjectMenu.Row(push),
		subjectMenu.Row(pop),
		subjectMenu.Row(letAhead),
	)

	return &Bot{
		b: b,

		startMenu: startMenu,
		editBtn:   &edit,
		chooseBtn: &choose,

		subjectMenu: subjectMenu,
		refreshBtn:  &refresh,
		pushBtn:     &push,
		popBtn:      &pop,
		letAheadBtn: &letAhead,

		cancel: cancel,
	}, ctx
}

func (b *Bot) StartMenu() *tele.ReplyMarkup {
	return b.startMenu
}

func (b *Bot) SubjectMenu() *tele.ReplyMarkup {
	return b.subjectMenu
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
		authorized.Handle(b.editBtn, handlers.Edit)
		authorized.Handle(b.chooseBtn, handlers.ChooseSubject)

		// Требует получить очередь из кеша
		authorized.Handle(b.refreshBtn, handlers.Refresh, middlewares.GetQueue)
		authorized.Handle(b.pushBtn, handlers.Push, middlewares.GetQueue)
		authorized.Handle(b.letAheadBtn, handlers.LetAhead, middlewares.GetQueue)

		// Нужны права админа
		authorized.Handle(b.popBtn, handlers.Pop, middlewares.GetQueue, middlewares.GetPermissions)
	}

	// Обработчик любого текста
	b.b.Handle(tele.OnText, handlers.OnText)
}
