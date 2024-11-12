package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"shortly/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

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
