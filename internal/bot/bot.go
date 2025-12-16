package bot

import (
	"context"

	"github.com/iskanye/mirea-queue/internal/config"
	"github.com/iskanye/mirea-queue/internal/interfaces"
	tele "gopkg.in/telebot.v4"
)

// Переменная для создания кнопок
var markup = tele.ReplyMarkup{}

type Bot struct {
	b *tele.Bot

	startMenu *tele.ReplyMarkup
	editBtn   *tele.Btn
	chooseBtn *tele.Btn

	subjectMenu *tele.ReplyMarkup
	returnBtn   *tele.Btn
	refreshBtn  *tele.Btn
	pushBtn     *tele.Btn
	letAheadBtn *tele.Btn
	popBtn      *tele.Btn
	clearBtn    *tele.Btn
	removeBtn   *tele.Btn

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

	// Кнопки
	editBtn := markup.Data("Изменить", "edit")
	chooseBtn := markup.Data("Очереди", "choose")
	returnBtn := markup.Data("Назад", "return")
	refreshBtn := markup.Data("Обновить", "update")
	pushBtn := markup.Data("Записаться", "push")
	popBtn := markup.Data("Позвать на сдачу", "pop")
	clearBtn := markup.Data("Очистить очередь", "clear")
	letAheadBtn := markup.Data("Пропустить в очереди", "let-ahead")
	removeBtn := markup.Data("Выйти из очереди", "remove")

	// Меню /start
	startMenu := &tele.ReplyMarkup{}
	startMenu.Inline(
		markup.Row(editBtn, chooseBtn),
	)

	// Меню предмета
	subjectMenu := &tele.ReplyMarkup{}
	subjectMenu.Inline(
		markup.Row(returnBtn, refreshBtn),
		markup.Row(pushBtn),
		markup.Row(letAheadBtn),
		markup.Row(removeBtn),
	)

	// Админ меню
	subjectAdminMenu := &tele.ReplyMarkup{}
	subjectAdminMenu.Inline(
		markup.Row(returnBtn, refreshBtn),
		markup.Row(pushBtn),
		markup.Row(letAheadBtn),
		markup.Row(removeBtn),
		markup.Row(popBtn),
		markup.Row(clearBtn),
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
		letAheadBtn: &letAheadBtn,
		popBtn:      &popBtn,
		clearBtn:    &clearBtn,
		removeBtn:   &removeBtn,

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
	// Автоматически ответь на колбеки
	b.b.Use(middlewares.CallbackRespond)

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
		authorized.Handle(b.removeBtn, handlers.Remove, middlewares.GetQueue)
		authorized.Handle(b.popBtn, handlers.Pop, middlewares.GetQueue)
		authorized.Handle(b.clearBtn, handlers.Clear, middlewares.GetQueue)
	}

	// Обработчик любого текста
	b.b.Handle(tele.OnText, handlers.OnText)

	// Кнопки выбора
	b.b.Handle("\f"+b.groupBtnUnique, handlers.ChooseGroup)
	b.b.Handle("\f"+b.subjectBtnUnique, handlers.ChooseSubjectButton)
}
