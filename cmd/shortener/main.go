package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"shortly/internal/app/api"
	"shortly/internal/app/config"
	"shortly/internal/app/errors"
	"shortly/internal/app/repository"
	"shortly/internal/app/service"
	"shortly/internal/compress"
	"shortly/internal/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	if err := run(ctx); err != nil {
		stop()
		log.Fatalf("Failed to start server: %v", err)
	}
}

func run(ctx context.Context) error {
	appLogger := logger.NewLogger()
	cfg := config.LoadConfig()
	repo, err := repository.NewFileStorageRepository(ctx, cfg.FileStoragePath)
	if err != nil {
		return errors.ErrFailedToInitializeRepository
	}
	router := setupRouter(cfg, *appLogger, repo)

	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.ListenAndServe()
	}()

	appLogger.Info().Msgf("Listening on %s", cfg.Addr)

	select {
	case <-ctx.Done():
		appLogger.Info().Msg("Shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err = server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown failed: %w", err)
		}

		appLogger.Info().Msg("Server gracefully stopped")

		repo.Wait()
		appLogger.Info().Msg("Repository shutdown completed")
		return nil
	case err = <-errCh:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
		return err
	}
}

func setupRouter(cfg *config.Config, appLogger logger.Logger, repo *repository.FileStorageRepository) http.Handler {
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
