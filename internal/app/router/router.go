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
	handler := api.NewURLHandler(cfg, shortener)

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

	router.Post("/api/shorten", handler.HandleCreateShortLink)
	router.Get("/api/shorten/{id}", handler.HandleGetShortLink)
	router.Post("/", handler.DeprecatedHandleCreateShortLink)
	router.Get("/{id}", handler.DeprecatedHandleGetShortLink)

	return router
}
