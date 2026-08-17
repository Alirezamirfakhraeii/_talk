package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/app"
	"github.com/ALirezamirfakhraeii/samatalk/backend/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf(
			"load configuration: %v",
			err,
		)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	application, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatalf("initialize application: %v", err)
	}
	defer application.Close()

	err = application.Run(ctx)
	if err != nil {
		log.Fatalf(
			"run application: %v",
			err,
		)
	}
}
