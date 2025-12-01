package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/iskanye/mirea-queue/internal/app"
	"github.com/iskanye/mirea-queue/internal/config"
)

func main() {
	cfg := config.MustLoadConfig()
	log := slog.Default()

	app := app.New(log, cfg)
	defer app.Stop()

	go func() {
		app.Run()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
}
