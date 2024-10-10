package app

import (
	"fmt"
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

func Run() {
	cfg := config.LoadConfig()

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
		middleware.Heartbeat("/health"),
		middleware.Logger,
		middleware.RequestID,
		middleware.Recoverer,
	)

	router.Post("/", handler.HandleCreateShortLink)
	router.Get("/{id}", handler.HandleGetShortLink)

	fmt.Println("Running server on", cfg.Addr)
	fmt.Println("Shortener base address is", cfg.BaseURL)
	fmt.Println("Shortener client address is", cfg.ClientURL)

	err := http.ListenAndServe(cfg.Addr, router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
