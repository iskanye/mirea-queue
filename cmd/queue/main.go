package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/iskanye/mirea-queue/internal/app"
)

func main() {
	app := app.New()

	go func() {
		app.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.Stop()
}
