package app

import (
	"context"
	"net/http"
	"time"

	"shortly/internal/app/config"
	"shortly/internal/app/repository"
	"shortly/internal/app/repository/persistence"
	"shortly/internal/app/router"
	"shortly/internal/app/server"
	"shortly/internal/app/worker"
	"shortly/internal/logger"
)

// Application is the main application structure
type Application struct {
	cfg                *config.Config
	logger             *logger.Logger
	persistenceManager persistence.Manager
	deleteWorker       worker.Worker
	server             server.Server
	pprofServer        server.PprofServer
}

// NewApplication creates a new application instance
func NewApplication(ctx context.Context) (*Application, error) {
	cfg := config.LoadConfig()
	appLogger := logger.NewLogger()

	appRepository, err := initRepository(ctx, cfg, appLogger)
	if err != nil {
		return nil, err
	}
	persistenceManager := persistence.NewPersistenceManager(cfg, appRepository, appLogger)

	deleteWorker := worker.NewDeleteWorker(ctx, cfg, appRepository, appLogger)
	deleteWorker.Start()

	appRouter := router.NewRouter(cfg, appRepository, deleteWorker, appLogger)
	appServer := server.NewServer(cfg, appRouter)
	pprofServer := server.NewPprofServer(cfg)

	return &Application{
		cfg:                cfg,
		logger:             appLogger,
		persistenceManager: persistenceManager,
		deleteWorker:       deleteWorker,
		server:             appServer,
		pprofServer:        pprofServer,
	}, nil
}

// Run starts the application
func (a *Application) Run(ctx context.Context) error {
	err := a.persistenceManager.Load()
	if err != nil {
		return err
	}

	serverErrors := make(chan error, 1)
	go func() {
		if err := a.server.Run(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	go func() {
		if err := a.pprofServer.Run(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	a.logger.Info().Msgf("Application starting in %s", a.cfg.AppEnv)
	a.logger.Info().Msgf("Listening on %s", a.cfg.Addr)
	a.logger.Info().Msgf("Profiler on %s", a.cfg.ProfilerAddr)

	select {
	case <-ctx.Done():
		a.logger.Info().Msg("Shutting down server...")

		a.deleteWorker.Stop()

		err = a.persistenceManager.Save()
		if err != nil {
			return err
		}

		shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err = a.server.Shutdown(shutdownCtx)
		if err != nil {
			return err
		}

		err = a.pprofServer.Shutdown(shutdownCtx)
		if err != nil {
			return err
		}

		a.logger.Info().Msg("Server gracefully stopped")
		return nil
	case err = <-serverErrors:
		return err
	}
}

func initRepository(ctx context.Context, cfg *config.Config, logger *logger.Logger) (repository.Repository, error) {
	repo, err := repository.NewRepository(ctx, &repository.Factory{
		DSN:    cfg.DatabaseDSN,
		Logger: logger,
	})

	if err != nil {
		logger.Error().Err(err).Msg("Failed to initialize application repository")
		return nil, err
	}

	return repo, nil
}
