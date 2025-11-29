package bot

import (
	"context"

	"github.com/iskanye/mirea-queue/internal/config"
	"github.com/iskanye/mirea-queue/internal/interfaces"
	tele "gopkg.in/telebot.v4"
)

type Bot struct {
	b *tele.Bot

	cancel context.CancelFunc
}

func New(
	cfg *config.Config,
) (*Bot, context.Context) {
	pref := tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: cfg.BotTimeout},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Bot{
		b:      b,
		cancel: cancel,
	}, ctx
}

func (b *Bot) Start() {
	b.b.Start()
}

func (b *Bot) Stop() {
	b.cancel()
	b.b.Stop()
}

func (b *Bot) RegisterHandlers(handlers interfaces.BotHandlers) {
	b.b.Handle("/start", handlers.Start)
}
