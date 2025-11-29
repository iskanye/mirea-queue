package app

import (
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/bot"
	"github.com/iskanye/mirea-queue/internal/config"
	botHandlers "github.com/iskanye/mirea-queue/internal/handlers/bot"
	"github.com/iskanye/mirea-queue/internal/repositories/postgres"
	"github.com/iskanye/mirea-queue/internal/repositories/redis"
	"github.com/iskanye/mirea-queue/internal/services/queue"
	"github.com/iskanye/mirea-queue/internal/services/users"
)

type App struct {
	log *slog.Logger
	bot *bot.Bot
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	redis, err := redis.New(cfg)
	if err != nil {
		panic(err)
	}

	postgres, err := postgres.New(cfg)
	if err != nil {
		panic(err)
	}

	queue := queue.New(log, redis, redis)
	users := users.New(log, postgres, postgres, postgres, postgres)

	handlers := botHandlers.New(log, queue, users)
	bot := bot.New(cfg, handlers)

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
