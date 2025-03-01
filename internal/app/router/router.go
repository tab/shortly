package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"shortly/internal/app/api"
	"shortly/internal/app/config"
	"shortly/internal/app/middleware/auth"
	"shortly/internal/app/middleware/compress"
	"shortly/internal/app/repository"
	"shortly/internal/app/service"
	"shortly/internal/app/worker"
	"shortly/internal/logger"
)

// NewRouter creates a new router instance
func NewRouter(cfg *config.Config, repo repository.Repository, worker worker.Worker, appLogger *logger.Logger) http.Handler {
	rand := service.NewSecureRandom()
	shortener := service.NewURLService(cfg, repo, rand, worker)
	shortenerHandler := api.NewURLHandler(cfg, shortener)

	health := service.NewHealthService(repo)
	healthHandler := api.NewHealthHandler(health)

	stats := service.NewStatsReporter(repo)
	statsHandler := api.NewStatsHandler(stats)

	authenticator := service.NewAuthService(cfg)

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

	router.Get("/live", healthHandler.HandleLiveness)
	router.Get("/ready", healthHandler.HandleReadiness)
	router.Get("/ping", healthHandler.HandlePing)

	router.Get("/api/internal/stats", statsHandler.HandleStats)

	// NOTE: protected routes
	router.Group(func(r chi.Router) {
		r.Use(auth.Middleware(authenticator))

		r.Get("/api/user/urls", shortenerHandler.HandleGetUserURLs)
		r.Delete("/api/user/urls", shortenerHandler.HandleBatchDeleteUserURLs)
	})

	// NOTE: public routes
	router.Group(func(r chi.Router) {
		r.Use(auth.Middleware(authenticator))
		r.Post("/api/shorten", shortenerHandler.HandleCreateShortLink)
		r.Get("/api/shorten/{id}", shortenerHandler.HandleGetShortLink)
		r.Post("/api/shorten/batch", shortenerHandler.HandleBatchCreateShortLink)
		r.Post("/", shortenerHandler.DeprecatedHandleCreateShortLink)
		r.Get("/{id}", shortenerHandler.DeprecatedHandleGetShortLink)
	})

	return router
}
