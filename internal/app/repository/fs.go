package repository

import (
	"encoding/json"
	"os"
	"sync"

	"shortly/internal/app/errors"
)

type FileStorageRepository struct {
	FilePath string
	sync.RWMutex
}

func NewFileStorageRepository(filePath string) *FileStorageRepository {
	return &FileStorageRepository{
		FilePath: filePath,
	}
}

func (store *FileStorageRepository) Set(url URL) error {
	store.Lock()
	defer store.Unlock()

	file, err := os.OpenFile(store.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.ErrFailedToOpenFile
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(url); err != nil {
		return errors.ErrFailedToWriteToFile
	}

	return nil
}

func (store *FileStorageRepository) Get(shortCode string) (*URL, bool) {
	store.RLock()
	defer store.RUnlock()

	file, err := os.OpenFile(store.FilePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, false
	}
	defer file.Close()

	results := json.NewDecoder(file)
	var url URL

	for results.More() {
		if err := results.Decode(&url); err != nil {
			return nil, false
		}

		if url.ShortCode == shortCode {
			return &url, true
		}
	}

	return nil, false
}
