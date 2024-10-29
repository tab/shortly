package persistence

import (
	"shortly/internal/app/repository"
	"shortly/internal/logger"
)

type Manager struct {
	repo      repository.Repository
	file      repository.FileRepository
	appLogger *logger.Logger
}

func NewPersistenceManager(repo repository.Repository, file repository.FileRepository, logger *logger.Logger) *Manager {
	return &Manager{
		repo:      repo,
		file:      file,
		appLogger: logger,
	}
}

func (pm *Manager) Load() error {
	snapshot, err := pm.file.Load()

	if err == nil {
		pm.repo.Restore(snapshot)
	} else {
		pm.appLogger.Error().Err(err).Msg("Failed to load data from file")
	}

	return nil
}

func (pm *Manager) Save() error {
	snapshot := pm.repo.CreateMemento()

	if err := pm.file.Save(snapshot); err != nil {
		pm.appLogger.Error().Err(err).Msg("Failed to save data to file")
	}

	return nil
}
