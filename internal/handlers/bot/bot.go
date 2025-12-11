package bot

import (
	"context"
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/interfaces"
	"gopkg.in/telebot.v4"
)

type Bot struct {
	log *slog.Logger
	ctx context.Context

	startMenu        *telebot.ReplyMarkup
	subjectMenu      *telebot.ReplyMarkup
	subjectAdminMenu *telebot.ReplyMarkup

	// id для кнопок выбора группы и предмета
	groupBtnUnique   string
	subjectBtnUnique string

	channels map[int64]chan string

	queueService    interfaces.QueueService
	usersService    interfaces.UsersService
	adminService    interfaces.AdminService
	scheduleService interfaces.ScheduleService
}

func New(
	log *slog.Logger,
	ctx context.Context,
	startMenu *telebot.ReplyMarkup,
	subjectMenu *telebot.ReplyMarkup,
	subjectAdminMenu *telebot.ReplyMarkup,
	groupBtnUnique string,
	subjectBtnUnique string,
	queueService interfaces.QueueService,
	usersService interfaces.UsersService,
	adminService interfaces.AdminService,
	scheduleService interfaces.ScheduleService,
) *Bot {
	return &Bot{
		log: log,
		ctx: ctx,

		startMenu:        startMenu,
		subjectMenu:      subjectMenu,
		subjectAdminMenu: subjectAdminMenu,

		channels: make(map[int64]chan string),

		groupBtnUnique:   groupBtnUnique,
		subjectBtnUnique: subjectBtnUnique,

		queueService:    queueService,
		usersService:    usersService,
		adminService:    adminService,
		scheduleService: scheduleService,
	}
}

// Обрабатывает текстовый ввод пользователя
func (b *Bot) OnText(c telebot.Context) error {
	if ch, ok := b.channels[c.Chat().ID]; ok {
		ch <- c.Message().Text

		err := c.Bot().Delete(c.Message())
		if err != nil {
			return err
		}
		return nil
	}

	b.log.Warn("Unhandled input",
		slog.String("text", c.Text()),
	)
	return nil
}

// Оборачивает команду в диалог, ожидая ввода от пользователя.
// Ввод пользователя передается через канал ch
func (b *Bot) Dialogue(
	c telebot.Context,
	fun func(ch <-chan string, c telebot.Context) error,
) error {
	chatID := c.Chat().ID
	ch := make(chan string, 1)
	b.channels[chatID] = ch

	defer func(cmd string) {
		close(ch)
		delete(b.channels, chatID)

		b.log.Info("Dialogue chain closed")
	}(c.Text())

	b.log.Info("Dialogue chain started")
	return fun(ch, c)
}
