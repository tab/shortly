package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"shortly/internal/app/api"
	"shortly/internal/app/config"
	"shortly/internal/app/repository"
	"shortly/internal/app/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func run() error {
	cfg := config.LoadConfig()
	router := setupRouter(cfg)

	return http.ListenAndServe(cfg.Addr, router)
}

func setupRouter(cfg *config.Config) http.Handler {
	repo := repository.NewInMemoryRepository()
	rand := service.NewSecureRandom()
	shortener := service.NewURLService(repo, rand, cfg)
	handler := api.NewURLHandler(cfg, shortener)

	router := chi.NewRouter()

	router.Use(
		cors.Handler(cors.Options{
			AllowedOrigins: []string{cfg.ClientURL},
			AllowedMethods: []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type"},
			MaxAge:         300,
		}),
		middleware.Logger,
		middleware.RequestID,
		middleware.Recoverer,
		middleware.Heartbeat("/health"),
	)

	router.Post("/", handler.HandleCreateShortLink)
	router.Get("/{id}", handler.HandleGetShortLink)

	return router
}
