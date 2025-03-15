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
	"shortly/internal/app/service"
	"shortly/internal/app/version"
	"shortly/internal/app/worker"
	"shortly/internal/logger"
)

const shutdownTimeout = 5 * time.Second

// Application is the main application structure
type Application struct {
	cfg                *config.Config
	logger             *logger.Logger
	shortener          *service.URLService
	persistenceManager persistence.Manager
	deleteWorker       worker.Worker
	server             server.Server
	grpcServer         server.GRPCServer
	grpcGateway        server.GRPCGateway
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

	rand := service.NewSecureRandom()
	shortener := service.NewURLService(cfg, appRepository, rand, deleteWorker)

	appRouter := router.NewRouter(cfg, shortener, appRepository, appLogger)
	appServer := server.NewServer(cfg, appRouter)

	grpcServer := server.NewGRPCServer(cfg, shortener)
	grpcGateway := server.NewGRPCGateway(ctx, cfg)
	pprofServer := server.NewPprofServer(cfg)

	return &Application{
		cfg:                cfg,
		logger:             appLogger,
		shortener:          shortener,
		persistenceManager: persistenceManager,
		deleteWorker:       deleteWorker,
		server:             appServer,
		grpcServer:         grpcServer,
		grpcGateway:        grpcGateway,
		pprofServer:        pprofServer,
	}, nil
}

// Run starts the application
func (a *Application) Run(ctx context.Context) error {
	if err := a.persistenceManager.Load(); err != nil {
		return err
	}

	serverErrors := make(chan error, 4)

	go func() {
		if err := a.server.Run(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	go func() {
		if err := a.grpcServer.Run(); err != nil {
			serverErrors <- err
		}
	}()

	go func() {
		if err := a.grpcGateway.Run(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	go func() {
		if err := a.pprofServer.Run(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	appVersion := version.NewVersion()

	a.logger.Info().Msgf("Build version: %s", appVersion.Version())
	a.logger.Info().Msgf("Build date: %s", appVersion.Date())
	a.logger.Info().Msgf("Build commit: %s", appVersion.Commit())

	a.logger.Info().Msgf("Application starting in %s", a.cfg.AppEnv)
	a.logger.Info().Msgf("Listening on %s", a.cfg.Addr)
	a.logger.Info().Msgf("gRPC listening on %s", a.cfg.GRPCServerAddr)
	a.logger.Info().Msgf("gRPC-Gateway listening on %s", a.cfg.GRPCGatewayAddr)
	a.logger.Info().Msgf("Profiler on %s", a.cfg.ProfilerAddr)

	select {
	case <-ctx.Done():
		a.logger.Info().Msg("Shutting down...")

		a.deleteWorker.Stop()

		if err := a.persistenceManager.Save(); err != nil {
			return err
		}

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		a.logger.Info().Msg("HTTP server gracefully stopped")

		if err := a.grpcServer.Shutdown(shutdownCtx); err != nil {
			return err
		}
		a.logger.Info().Msg("gRPC server gracefully stopped")

		if err := a.grpcGateway.Shutdown(shutdownCtx); err != nil {
			return err
		}
		a.logger.Info().Msg("gRPC-Gateway server gracefully stopped")

		if err := a.pprofServer.Shutdown(shutdownCtx); err != nil {
			return err
		}
		a.logger.Info().Msg("Profiler server gracefully stopped")

		return nil
	case err := <-serverErrors:
		return err
	}
}

// initRepository initializes the application repository based on the configuration
func initRepository(ctx context.Context, cfg *config.Config, logger *logger.Logger) (repository.Repository, error) {
	if cfg.DatabaseDSN == "" {
		return repository.NewInMemoryRepository(), nil
	}

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
