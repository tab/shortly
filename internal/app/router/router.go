package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"shortly/internal/app/api"
	"shortly/internal/app/config"
	"shortly/internal/app/repository"
	"shortly/internal/app/service"
	"shortly/internal/compress"
	"shortly/internal/logger"
)

func NewRouter(cfg *config.Config, appLogger *logger.Logger, repo repository.Repository) http.Handler {
	rand := service.NewSecureRandom()
	shortener := service.NewURLService(cfg, repo, rand)
	shortenerHandler := api.NewURLHandler(cfg, shortener)

	health := service.NewHealthService(repo)
	healthHandler := api.NewHealthHandler(health)

	router := chi.NewRouter()
	router.Use(
		cors.Handler(cors.Options{
			AllowedOrigins: []string{cfg.ClientURL},
			AllowedMethods: []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type"},
			MaxAge:         300,
		}),
		appLogger.Middleware,
		compress.Middleware,
	)

	router.Get("/ping", healthHandler.HandlePing)
	router.Post("/api/shorten", shortenerHandler.HandleCreateShortLink)
	router.Get("/api/shorten/{id}", shortenerHandler.HandleGetShortLink)
	router.Post("/api/shorten/batch", shortenerHandler.HandleBatchCreateShortLink)
	router.Post("/", shortenerHandler.DeprecatedHandleCreateShortLink)
	router.Get("/{id}", shortenerHandler.DeprecatedHandleGetShortLink)

	return router
}
