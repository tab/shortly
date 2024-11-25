package app

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortly/internal/app/config"
	"shortly/internal/app/repository"
	"shortly/internal/app/repository/persistence"
	"shortly/internal/app/server"
	"shortly/internal/app/worker"
	"shortly/internal/logger"
)

func Test_NewApplication(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name     string
		expected Application
	}{
		{
			name:     "Success",
			expected: Application{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := NewApplication(ctx)
			assert.NoError(t, err)

			assert.NotNil(t, app)
			assert.NotNil(t, app.cfg)
			assert.NotNil(t, app.logger)
			assert.NotNil(t, app.server)
		})
	}
}

func Test_initRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	appLogger := logger.NewLogger()

	type result struct {
		repoType interface{}
		err      error
	}

	tests := []struct {
		name     string
		dsn      string
		expected result
	}{
		{
			name: "With valid DSN",
			dsn:  "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
			expected: result{
				repoType: &repository.DatabaseRepo{},
				err:      nil,
			},
		},
		{
			name: "Without DSN",
			dsn:  "",
			expected: result{
				repoType: &repository.InMemoryRepo{},
				err:      nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				DatabaseDSN: tt.dsn,
			}

			repo, err := initRepository(ctx, cfg, appLogger)

			if tt.expected.err != nil {
				assert.Error(t, err)
				assert.Nil(t, repo)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, repo)
				assert.IsType(t, tt.expected.repoType, repo)
			}
		})
	}
}

func Test_Application_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPersistenceManager := persistence.NewMockManager(ctrl)
	mockServer := server.NewMockServer(ctrl)
	appLogger := logger.NewLogger()
	repo := repository.NewInMemoryRepository()
	appWorker := worker.NewDeleteWorker(&config.Config{}, repo, appLogger)

	tests := []struct {
		name     string
		before   func()
		expected error
	}{
		{
			name: "Success",
			before: func() {
				mockPersistenceManager.EXPECT().Load().Return(nil)
				mockServer.EXPECT().Run().DoAndReturn(func() error {
					time.Sleep(100 * time.Millisecond)
					return nil
				})
				mockPersistenceManager.EXPECT().Save().Return(nil)
				mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil)
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			app := &Application{
				cfg:                &config.Config{},
				logger:             appLogger,
				persistenceManager: mockPersistenceManager,
				deleteWorker:       appWorker,
				server:             mockServer,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer cancel()

			err := app.Run(ctx)

			if tt.expected != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
