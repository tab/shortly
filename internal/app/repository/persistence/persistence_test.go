package persistence

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortly/internal/app/config"
	"shortly/internal/app/errors"
	"shortly/internal/app/repository"
	"shortly/internal/logger"
	"shortly/internal/spec"
)

func TestMain(m *testing.M) {
	if err := spec.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	if os.Getenv("GO_ENV") == "ci" {
		os.Exit(0)
	}

	code := m.Run()
	os.Exit(code)
}

func Test_NewPersistenceManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filePath := t.TempDir() + "/store-test.json"
	appLogger := logger.NewLogger()

	ctx := context.Background()
	dsn := os.Getenv("DATABASE_DSN")
	databaseRepo, _ := repository.NewDatabaseRepository(ctx, dsn)
	inMemoryRepo := repository.NewInMemoryRepository()
	fileRepo := repository.NewFileRepository(filePath)

	tests := []struct {
		name         string
		cfg          *config.Config
		repo         repository.Repository
		expectedType interface{}
	}{
		{
			name: "InMemory repo with file path",
			cfg: &config.Config{
				FileStoragePath: filePath,
			},
			repo: inMemoryRepo,
			expectedType: &manager{
				repo:      inMemoryRepo,
				file:      fileRepo,
				appLogger: appLogger,
			},
		},
		{
			name: "InMemory repo without file path",
			cfg: &config.Config{
				FileStoragePath: "",
			},
			repo:         inMemoryRepo,
			expectedType: &noOpManager{},
		},
		{
			name: "Database with file path",
			cfg: &config.Config{
				FileStoragePath: filePath,
			},
			repo:         databaseRepo,
			expectedType: &noOpManager{},
		},
		{
			name: "Database without file path",
			cfg: &config.Config{
				FileStoragePath: "",
			},
			repo:         databaseRepo,
			expectedType: &noOpManager{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pm := NewPersistenceManager(tt.cfg, appLogger, tt.repo)
			assert.NotNil(t, pm)
			assert.IsType(t, tt.expectedType, pm)
		})
	}
}

func Test_noOpManager_Load(t *testing.T) {
	n := &noOpManager{}
	err := n.Load()
	assert.NoError(t, err)
}

func Test_noOpManager_Save(t *testing.T) {
	n := &noOpManager{}
	err := n.Save()
	assert.NoError(t, err)
}

func Test_PersistenceManager_Load(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filePath := t.TempDir() + "/store-test.json"
	cfg := &config.Config{
		FileStoragePath: filePath,
	}
	mockRepo := repository.NewMockInMemory(ctrl)
	mockFileRepo := repository.NewMockFile(ctrl)
	appLogger := logger.NewLogger()

	UUID, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f726639")

	snapshot := &repository.Memento{
		State: []repository.URL{
			{
				UUID:      UUID,
				LongURL:   "http://example.com",
				ShortCode: "abcd1234",
			},
		},
	}

	tests := []struct {
		name   string
		before func()
	}{
		{
			name: "Success",
			before: func() {
				mockFileRepo.EXPECT().Load().Return(snapshot, nil)
				mockRepo.EXPECT().Restore(snapshot)
			},
		},
		{
			name: "Failure",
			before: func() {
				mockFileRepo.EXPECT().Load().Return(nil, errors.ErrFailedToOpenFile)
				appLogger.Error().Err(errors.ErrFailedToOpenFile).Msg("Failed to load data from file")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			pm := &manager{
				repo:      mockRepo,
				file:      mockFileRepo,
				appLogger: appLogger,
			}
			err := pm.Load()

			assert.NoError(t, err)

			t.Cleanup(func() {
				os.RemoveAll(cfg.FileStoragePath)
			})
		})
	}
}

func Test_PersistenceManager_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filePath := t.TempDir() + "/store-test.json"
	cfg := &config.Config{
		FileStoragePath: filePath,
	}
	mockRepo := repository.NewMockInMemory(ctrl)
	mockFileRepo := repository.NewMockFile(ctrl)
	appLogger := logger.NewLogger()

	UUID, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f726639")

	snapshot := &repository.Memento{
		State: []repository.URL{
			{
				UUID:      UUID,
				LongURL:   "http://example.com",
				ShortCode: "abcd1234",
			},
		},
	}

	tests := []struct {
		name   string
		before func()
	}{
		{
			name: "Success",
			before: func() {
				mockRepo.EXPECT().CreateMemento().Return(snapshot)
				mockFileRepo.EXPECT().Save(snapshot).Return(nil)
			},
		},
		{
			name: "Failure",
			before: func() {
				mockRepo.EXPECT().CreateMemento().Return(snapshot)
				mockFileRepo.EXPECT().Save(snapshot).Return(errors.ErrFailedToOpenFile)
				appLogger.Error().Err(errors.ErrFailedToOpenFile).Msg("Failed to save data to file")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			pm := &manager{
				repo:      mockRepo,
				file:      mockFileRepo,
				appLogger: appLogger,
			}
			err := pm.Save()

			assert.NoError(t, err)

			t.Cleanup(func() {
				os.RemoveAll(cfg.FileStoragePath)
			})
		})
	}
}
