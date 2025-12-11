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
	returnBtn   *tele.Btn
	refreshBtn  *tele.Btn
	pushBtn     *tele.Btn
	popBtn      *tele.Btn
	letAheadBtn *tele.Btn

	// id для кнопок выбора группы и предмета
	groupBtnUnique   string
	subjectBtnUnique string

	subjectAdminMenu *tele.ReplyMarkup

	cancel context.CancelFunc
}

func New(
	cfg *config.Config,
	groupBtnUnique string,
	subjectBtnUnique string,
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

	// Меню /start
	startMenu := &tele.ReplyMarkup{}
	editBtn := startMenu.Data("Изменить", "edit")
	chooseBtn := startMenu.Data("Очереди", "choose")
	startMenu.Inline(
		startMenu.Row(editBtn, chooseBtn),
	)

	// Меню предмета
	subjectMenu := &tele.ReplyMarkup{}
	returnBtn := subjectMenu.Data("Назад", "return")
	refreshBtn := subjectMenu.Data("Обновить", "update")
	pushBtn := subjectMenu.Data("Записаться", "push")
	popBtn := subjectMenu.Data("Позвать на сдачу", "pop")
	letAheadBtn := subjectMenu.Data("Пропустить в очереди", "let-ahead")
	subjectMenu.Inline(
		subjectMenu.Row(returnBtn, refreshBtn),
		subjectMenu.Row(pushBtn),
		subjectMenu.Row(letAheadBtn),
	)

	// Админ меню
	subjectAdminMenu := &tele.ReplyMarkup{}
	subjectAdminMenu.Inline(
		subjectAdminMenu.Row(returnBtn, refreshBtn),
		subjectAdminMenu.Row(pushBtn),
		subjectAdminMenu.Row(popBtn),
		subjectAdminMenu.Row(letAheadBtn),
	)

	return &Bot{
		b: b,

		startMenu: startMenu,
		editBtn:   &editBtn,
		chooseBtn: &chooseBtn,

		subjectMenu: subjectMenu,
		returnBtn:   &returnBtn,
		refreshBtn:  &refreshBtn,
		pushBtn:     &pushBtn,
		popBtn:      &popBtn,
		letAheadBtn: &letAheadBtn,

		subjectAdminMenu: subjectAdminMenu,

		groupBtnUnique:   groupBtnUnique,
		subjectBtnUnique: subjectBtnUnique,

		cancel: cancel,
	}, ctx
}

func (b *Bot) StartMenu() *tele.ReplyMarkup {
	return b.startMenu
}

func (b *Bot) SubjectMenu() *tele.ReplyMarkup {
	return b.subjectMenu
}

func (b *Bot) SubjectAdminMenu() *tele.ReplyMarkup {
	return b.subjectAdminMenu
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
		authorized.Handle(b.returnBtn, handlers.Return)

		// Требует получить очередь из кеша
		authorized.Handle(b.refreshBtn, handlers.Refresh, middlewares.GetQueue)
		authorized.Handle(b.pushBtn, handlers.Push, middlewares.GetQueue)
		authorized.Handle(b.letAheadBtn, handlers.LetAhead, middlewares.GetQueue)

		// Нужны права админа
		authorized.Handle(b.popBtn, handlers.Pop, middlewares.GetQueue, middlewares.GetPermissions)
	}

	// Обработчик любого текста
	b.b.Handle(tele.OnText, handlers.OnText)

	// Кнопки выбора
	b.b.Handle("\f"+b.groupBtnUnique, handlers.ChooseGroup)
}
