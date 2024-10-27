package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"shortly/internal/app/errors"
)

const (
	FileDumpInterval = 15 * time.Second
)

type FileStorageRepository struct {
	FilePath string
	InMemory *InMemoryRepository
	wg       sync.WaitGroup
}

func NewFileStorageRepository(ctx context.Context, filePath string) (*FileStorageRepository, error) {
	repo := &FileStorageRepository{
		FilePath: filePath,
		InMemory: NewInMemoryRepository(),
	}

	if err := repo.loadFromFile(); err != nil {
		return nil, fmt.Errorf("failed to load data from file: %w", err)
	}

	repo.wg.Add(1)
	go repo.snapshot(ctx)

	return repo, nil
}

func (store *FileStorageRepository) Set(url URL) error {
	if err := store.InMemory.Set(url); err != nil {
		return err
	}

	return nil
}

func (store *FileStorageRepository) Get(shortCode string) (*URL, bool) {
	return store.InMemory.Get(shortCode)
}

func (store *FileStorageRepository) GetAll() []URL {
	return store.InMemory.GetAll()
}

func (store *FileStorageRepository) Load() error {
	return store.loadFromFile()
}

func (store *FileStorageRepository) Save() error {
	return store.saveToFile()
}

func (store *FileStorageRepository) Wait() {
	store.wg.Wait()
}

func (store *FileStorageRepository) loadFromFile() error {
	fileInfo, err := os.Stat(store.FilePath)
	if err == nil && fileInfo.IsDir() {
		return errors.ErrFilePathIsDirectory
	}

	file, err := os.OpenFile(store.FilePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return errors.ErrFailedToOpenFile
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	for {
		var url URL
		if err = decoder.Decode(&url); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		err = store.InMemory.Set(url)
		if err != nil {
			return err
		}
	}

	return nil
}

func (store *FileStorageRepository) saveToFile() error {
	fileInfo, err := os.Stat(store.FilePath)
	if err == nil && fileInfo.IsDir() {
		return errors.ErrFilePathIsDirectory
	}

	data := store.InMemory.GetAll()

	file, err := os.OpenFile(store.FilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.ErrFailedToOpenFile
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for _, url := range data {
		if err = encoder.Encode(url); err != nil {
			return errors.ErrFailedToWriteToFile
		}
	}

	return nil
}

func (store *FileStorageRepository) snapshot(ctx context.Context) {
	defer store.wg.Done()

	ticker := time.NewTicker(FileDumpInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := store.saveToFile(); err != nil {
				log.Printf("Failed to write to file: %v", err)
			}
		case <-ctx.Done():
			if err := store.saveToFile(); err != nil {
				log.Printf("Failed to write to file during shutdown: %v", err)
			}
			return
		}
	}
}
