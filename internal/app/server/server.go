package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"shortly/internal/app/config"
)

var options *config.AppConfig

func AppRouter() chi.Router {
	router := chi.NewRouter()

	router.Use(
		cors.Handler(cors.Options{
			AllowedOrigins: []string{options.Flags.ClientURL},
			AllowedMethods: []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type"},
			MaxAge:         300,
		}),
		middleware.Heartbeat("/status"),
		middleware.Logger,
		middleware.RequestID,
		middleware.Recoverer)

	router.Post("/", HandleCreateShortLink)
	router.Get("/{id}", HandleGetShortLink)

	return router
}

func Run() {
	options = config

	fmt.Println("Running server on", options.Addr)
	fmt.Println("Shortener base address is", options.BaseURL)
	fmt.Println("Shortener client address is", options.ClientURL)

	err := http.ListenAndServe(options.Addr, AppRouter())
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
