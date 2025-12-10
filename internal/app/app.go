package app

import (
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/bot"
	"github.com/iskanye/mirea-queue/internal/config"
	botHandlers "github.com/iskanye/mirea-queue/internal/handlers/bot"
	botMiddlewares "github.com/iskanye/mirea-queue/internal/middlewares/bot"
	"github.com/iskanye/mirea-queue/internal/repositories/postgres"
	"github.com/iskanye/mirea-queue/internal/repositories/redis"
	"github.com/iskanye/mirea-queue/internal/services/admin"
	"github.com/iskanye/mirea-queue/internal/services/queue"
	"github.com/iskanye/mirea-queue/internal/services/users"
)

// Пагинация очереди
const queueRange = 10

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

	queue := queue.New(log, queueRange, redis, redis, redis, redis, redis, redis)
	users := users.New(log, postgres, postgres, postgres, postgres)
	admin := admin.New(log, cfg)

	bot, ctx := bot.New(cfg)
	handlers := botHandlers.New(log, ctx,
		bot.StartMenu(),
		bot.SubjectMenu(),
		bot.SubjectAdminMenu(),
		queue, users, admin,
	)
	middlewares := botMiddlewares.New(log, ctx, queue, users, admin)

	bot.Register(handlers, middlewares)

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
