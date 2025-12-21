package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iskanye/mirea-queue/internal/config"
	"github.com/iskanye/mirea-queue/internal/repositories/redis"
	"github.com/robfig/cron/v3"
)

func main() {
	// Загружаем московский временной пояс
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		panic(err)
	}

	c := cron.New(cron.WithLocation(location))

	cfg := config.MustLoadConfig()
	log := slog.Default()

	redis, err := redis.New(cfg)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Каждый раз в 6 часов по МСК очищает базу данных очереди
	_, err = c.AddFunc("0 6 * * *", func() {
		err := redis.FlushDB(ctx)
		if err != nil {
			log.Error("UNKNOWN ERROR",
				slog.String("err", err.Error()),
			)
		}
	})
	if err != nil {
		panic(err)
	}

	// Стартуем крон
	c.Start()
	defer c.Stop()

	log.Info("Cron started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	log.Info("Stopped")
}
