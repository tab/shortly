package routers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"shortly/internal/app/api/handlers"
	"shortly/internal/app/config"
	"shortly/internal/app/helpers"
	"shortly/internal/app/store"
)

func AppRouter(appConfig *config.AppConfig) (chi.Router, *handlers.Handler) {
	router := chi.NewRouter()

	handler := &handlers.Handler{
		AppConfig:    appConfig,
		SecureRandom: helpers.NewSecureRandom(),
		Store:        *store.NewURLStore(),
	}

	router.Use(
		cors.Handler(cors.Options{
			AllowedOrigins: []string{appConfig.ClientURL},
			AllowedMethods: []string{"GET", "POST", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type"},
			MaxAge:         300,
		}),
		middleware.Heartbeat("/status"),
		middleware.Logger,
		middleware.RequestID,
		middleware.Recoverer)

	router.Post("/", handler.HandleCreateShortLink)
	router.Get("/{id}", handler.HandleGetShortLink)

	return router, handler
}

func Run(appConfig *config.AppConfig) {
	fmt.Println("Running server on", appConfig.Addr)
	fmt.Println("Shortener base address is", appConfig.BaseURL)
	fmt.Println("Shortener client address is", appConfig.ClientURL)

	router, _ := AppRouter(appConfig)

	err := http.ListenAndServe(appConfig.Addr, router)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
