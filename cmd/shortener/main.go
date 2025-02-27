package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"shortly/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	application, err := app.NewApplication(ctx)
	if err != nil {
		stop()
		log.Fatalf("Failed to initialize application: %v", err)
	}

	if err = application.Run(ctx); err != nil {
		stop()
		log.Fatalf("Application error: %v", err)
	}
}
