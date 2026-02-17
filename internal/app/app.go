package app

import (
	"log/slog"

	"github.com/iskanye/mirea-queue/internal/bot"
	scheduleClient "github.com/iskanye/mirea-queue/internal/client/schedule"
	"github.com/iskanye/mirea-queue/internal/config"
	botHandlers "github.com/iskanye/mirea-queue/internal/handlers/bot"
	botMiddlewares "github.com/iskanye/mirea-queue/internal/middlewares/bot"
	"github.com/iskanye/mirea-queue/internal/repositories/postgres"
	"github.com/iskanye/mirea-queue/internal/repositories/redis"
	"github.com/iskanye/mirea-queue/internal/services/admin"
	"github.com/iskanye/mirea-queue/internal/services/queue"
	"github.com/iskanye/mirea-queue/internal/services/schedule"
	"github.com/iskanye/mirea-queue/internal/services/users"
)

const (
	// Пагинация очереди
	queueRange = 10
	// Пагинация групп
	groupRange = 5

	// id для кнопок выбора группы и предмета
	groupBtnUnique   = "group"
	subjectBtnUnique = "subject"
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

	client := scheduleClient.New()

	queue := queue.New(log, queueRange, redis, redis, redis, redis, redis, redis, redis)
	users := users.New(log, postgres, postgres, postgres, postgres)
	admin := admin.New(log, cfg)
	schedule := schedule.New(log, groupRange, client, client)

	bot, ctx := bot.New(log, cfg, groupBtnUnique, subjectBtnUnique)
	handlers := botHandlers.New(log, ctx,
		bot.StartMenu(),
		bot.SubjectMenu(),
		bot.SubjectAdminMenu(),
		groupBtnUnique, subjectBtnUnique,
		queue, users, admin, schedule,
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
