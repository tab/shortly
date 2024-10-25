package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"shortly/internal/app/api"
	"shortly/internal/app/config"
	"shortly/internal/app/repository"
	"shortly/internal/app/service"
	"shortly/internal/compress"
	"shortly/internal/logger"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func run() error {
	cfg := config.LoadConfig()
	router := setupRouter(cfg)

	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return server.ListenAndServe()
}

func setupRouter(cfg *config.Config) http.Handler {
	appLogger := logger.NewLogger()
	repo := repository.NewFileStorageRepository(cfg.FileStoragePath)
	rand := service.NewSecureRandom()
	shortener := service.NewURLService(cfg, repo, rand)
	handler := api.NewURLHandler(cfg, shortener)

	router := chi.NewRouter()

	router.Use(
		cors.Handler(cors.Options{
			AllowedOrigins: []string{cfg.ClientURL},
			AllowedMethods: []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type"},
			MaxAge:         300,
		}),
		// middleware.Logger,
		appLogger.Middleware,
		compress.Middleware,
		middleware.RequestID,
		middleware.Recoverer,
		middleware.Heartbeat("/health"),
	)

	router.Post("/api/shorten", handler.HandleCreateShortLink)
	router.Get("/api/shorten/{id}", handler.HandleGetShortLink)

	router.Post("/", handler.DeprecatedHandleCreateShortLink)
	router.Get("/{id}", handler.DeprecatedHandleGetShortLink)

	return router
}
