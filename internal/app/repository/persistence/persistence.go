package persistence

import (
	"shortly/internal/app/config"
	"shortly/internal/app/repository"
	"shortly/internal/logger"
)

// Manager is an interface for the persistence manager
type Manager interface {
	Load() error
	Save() error
}

// noOpManager is a no-op persistence manager
type noOpManager struct{}

// manager is a persistence manager
type manager struct {
	repo      repository.InMemory
	file      repository.File
	appLogger *logger.Logger
}

// NewPersistenceManager creates a new persistence manager instance
func NewPersistenceManager(cfg *config.Config, repo repository.Repository, logger *logger.Logger) Manager {
	inMemoryRepo, ok := repo.(repository.InMemory)
	if !ok {
		logger.Warn().Msg("Persistence manager initialization skipped, not an in-memory repository")
		return &noOpManager{}
	}

	if cfg.FileStoragePath == "" {
		logger.Warn().Msg("Persistence manager initialization skipped, file is not set")
		return &noOpManager{}
	}

	fileRepo := repository.NewFileRepository(cfg.FileStoragePath)
	logger.Info().Msg("Persistence manager is initialized with " + cfg.FileStoragePath)

	return &manager{
		repo:      inMemoryRepo,
		file:      fileRepo,
		appLogger: logger,
	}
}

// Load stub implementation for no-op manager
func (n *noOpManager) Load() error {
	return nil
}

// Load loads data from file and restores the repository state
func (pm *manager) Load() error {
	snapshot, err := pm.file.Load()

	if err == nil {
		pm.repo.Restore(snapshot)
	} else {
		pm.appLogger.Error().Err(err).Msg("Failed to load data from file")
	}

	return nil
}

// Save stub implementation for no-op manager
func (n *noOpManager) Save() error {
	return nil
}

// Save saves the repository state to the file
func (pm *manager) Save() error {
	snapshot := pm.repo.CreateMemento()

	if err := pm.file.Save(snapshot); err != nil {
		pm.appLogger.Error().Err(err).Msg("Failed to save data to file")
	}

	return nil
}
