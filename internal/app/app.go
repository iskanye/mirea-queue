package app

import (
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/bot"
	"github.com/iskanye/mirea-queue/internal/config"
)

type App struct {
	log *slog.Logger
	bot *bot.Bot
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	bot := bot.New(cfg)

	return &App{
		log: log,
		bot: bot,
	}
}

func (a *App) Run() {
	a.log.Info("Bot started")
	a.bot.Start()
}

func (a *App) Stop() {
	a.bot.Stop()
	a.log.Info("Bot stopped")
}
