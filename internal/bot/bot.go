package bot

import (
	"github.com/iskanye/mirea-queue/internal/config"
	"github.com/iskanye/mirea-queue/internal/interfaces"
	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	b *tele.Bot

	handlers interfaces.BotHandlers
}

func New(
	cfg *config.Config,
	handlers interfaces.BotHandlers,
) *Bot {
	pref := tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: cfg.BotTimeout},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		panic(err)
	}

	return &Bot{
		b:        b,
		handlers: handlers,
	}
}

func (b *Bot) Start() {
	b.registerHandlers()
	b.b.Start()
}

func (b *Bot) Stop() {
	b.b.Stop()
}

func (b *Bot) registerHandlers() {
	b.b.Handle("/start", b.handlers.Start)
}
