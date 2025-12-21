package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/iskanye/mirea-queue/internal/config"
	"github.com/iskanye/mirea-queue/internal/repositories/redis"
	"github.com/robfig/cron/v3"
)

func main() {
	const op = "cron"

	// Загружаем cron
	c := cron.New()

	cfg := config.MustLoadConfig()
	log := slog.With(
		slog.String("op", op),
	)

	redis, err := redis.New(cfg)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Очищает базу данных очереди по заданному правилу
	_, err = c.AddFunc(cfg.CronTab, func() {
		err := redis.FlushDB(ctx)
		if err != nil {
			log.Error("Unknown error",
				slog.String("err", err.Error()),
			)
		} else {
			log.Info("Flushed database")
		}
	})
	if err != nil {
		panic(err)
	}

	// Стартуем крон
	c.Start()
	defer c.Stop()

	log.Info("Started",
		slog.String("crontab", cfg.CronTab),
	)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	log.Info("Stopped")
}
