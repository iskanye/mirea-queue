package bot

import (
	"github.com/iskanye/mirea-queue/internal/config"
	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	b *tele.Bot
}

func New(cfg *config.Config) *Bot {
	pref := tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: cfg.BotTimeout},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		panic(err)
	}

	return &Bot{
		b: b,
	}
}

func (b *Bot) Start() {
	b.registerHandlers()
	b.b.Start()
}

func (b *Bot) Stop() {
	b.b.Stop()
}
