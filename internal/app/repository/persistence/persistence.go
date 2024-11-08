package persistence

import (
	"shortly/internal/app/config"
	"shortly/internal/app/repository"
	"shortly/internal/logger"
)

type Manager interface {
	Load() error
	Save() error
}

type noOpManager struct{}

type manager struct {
	repo      repository.InMemory
	file      repository.File
	appLogger *logger.Logger
}

func NewPersistenceManager(cfg *config.Config, logger *logger.Logger, repo repository.Repository) Manager {
	noOp := &noOpManager{}

	inMemoryRepo, ok := repo.(repository.InMemory)
	if !ok {
		logger.Warn().Msg("Persistence manager initialization skipped, not an in-memory repository")
		return noOp
	}

	if cfg.FileStoragePath == "" {
		logger.Warn().Msg("Persistence manager initialization skipped, file is not set")
		return noOp
	}

	fileRepo := repository.NewFileRepository(cfg.FileStoragePath)
	logger.Info().Msg("Persistence manager is initialized with " + cfg.FileStoragePath)

	return &manager{
		repo:      inMemoryRepo,
		file:      fileRepo,
		appLogger: logger,
	}
}

func (n *noOpManager) Load() error {
	return nil
}

func (pm *manager) Load() error {
	snapshot, err := pm.file.Load()

	if err == nil {
		pm.repo.Restore(snapshot)
	} else {
		pm.appLogger.Error().Err(err).Msg("Failed to load data from file")
	}

	return nil
}

func (n *noOpManager) Save() error {
	return nil
}

func (pm *manager) Save() error {
	snapshot := pm.repo.CreateMemento()

	if err := pm.file.Save(snapshot); err != nil {
		pm.appLogger.Error().Err(err).Msg("Failed to save data to file")
	}

	return nil
}
