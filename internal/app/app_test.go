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
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	tests := []struct {
		name  string
		error bool
	}{
		{
			name:  "Success",
			error: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := NewApplication(ctx)

			if tt.error {
				assert.Nil(t, app)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, app)
				assert.NotNil(t, app.cfg)
				assert.NotNil(t, app.logger)
				assert.NotNil(t, app.server)
				assert.NotNil(t, app.pprofServer)
			}
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

	cfg := &config.Config{
		Addr:      "localhost:8080",
		BaseURL:   "http://localhost:8080",
		ClientURL: "http://localhost:8080",
	}

	ctx := context.Background()
	mockPersistenceManager := persistence.NewMockManager(ctrl)
	mockServer := server.NewMockServer(ctrl)
	appLogger := logger.NewLogger()
	repo := repository.NewInMemoryRepository()
	appWorker := worker.NewDeleteWorker(ctx, cfg, repo, appLogger)
	mockPprofServer := server.NewMockPprofServer(ctrl)

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
				mockPprofServer.EXPECT().Run().DoAndReturn(func() error {
					time.Sleep(100 * time.Millisecond)
					return nil
				})
				mockPersistenceManager.EXPECT().Save().Return(nil)
				mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil)
				mockPprofServer.EXPECT().Shutdown(gomock.Any()).Return(nil)
			},
			expected: nil,
		},
		{
			name: "Error on Load",
			before: func() {
				mockPersistenceManager.EXPECT().Load().Return(assert.AnError)
			},
			expected: assert.AnError,
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
				pprofServer:        mockPprofServer,
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

func Test_Application_Run_ShutdownErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	cfg := &config.Config{
		Addr: "localhost:8080",
	}
	repo := repository.NewInMemoryRepository()
	appLogger := logger.NewLogger()
	appWorker := worker.NewDeleteWorker(ctx, cfg, repo, appLogger)

	mockPersistenceManager := persistence.NewMockManager(ctrl)
	mockServer := server.NewMockServer(ctrl)
	mockPprofServer := server.NewMockPprofServer(ctrl)

	mockPersistenceManager.EXPECT().Load().Return(nil).AnyTimes()
	mockServer.EXPECT().Run().DoAndReturn(func() error {
		time.Sleep(50 * time.Millisecond)
		return nil
	}).AnyTimes()
	mockPprofServer.EXPECT().Run().DoAndReturn(func() error {
		time.Sleep(50 * time.Millisecond)
		return nil
	}).AnyTimes()

	application := &Application{
		cfg:                cfg,
		logger:             appLogger,
		persistenceManager: mockPersistenceManager,
		deleteWorker:       appWorker,
		server:             mockServer,
		pprofServer:        mockPprofServer,
	}

	runErrCh := make(chan error, 1)

	tests := []struct {
		name     string
		before   func()
		expected error
	}{
		{
			name: "SaveError",
			before: func() {
				mockPersistenceManager.EXPECT().Save().Return(assert.AnError)
			},
			expected: assert.AnError,
		},
		{
			name: "ServerShutdownError",
			before: func() {
				mockPersistenceManager.EXPECT().Save().Return(nil)
				mockServer.EXPECT().Shutdown(gomock.Any()).Return(assert.AnError)
			},
			expected: assert.AnError,
		},
		{
			name: "PprofShutdownError",
			before: func() {
				mockPersistenceManager.EXPECT().Save().Return(nil)
				mockServer.EXPECT().Shutdown(gomock.Any()).Return(nil)
				mockPprofServer.EXPECT().Shutdown(gomock.Any()).Return(assert.AnError)
			},
			expected: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			go func() {
				err := application.Run(ctx)
				runErrCh <- err
			}()

			time.Sleep(50 * time.Millisecond)
			cancel()

			err := <-runErrCh

			if tt.expected != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
