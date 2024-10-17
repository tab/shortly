package main

import (
	"log"

	"shortly/internal/app"
)

func main() {
	err := app.Run(&app.HTTPServer{})

	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
