package repository

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortly/internal/app/config"
	"shortly/internal/logger"
)

func Test_NewRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	cfg := &config.Config{
		ClientURL:   config.ClientURL,
		DatabaseDSN: config.DatabaseDSN,
	}
	appLogger := logger.NewLogger()
	repo := NewMockRepository(ctrl)

	tests := []struct {
		name   string
		before func()
	}{
		{
			name: "PostgreSQL database repository",
			before: func() {
				repo.EXPECT().Ping(ctx).Return(nil)
			},
		},
		{
			name: "In-Memory repository",
			before: func() {
				repo.EXPECT().Ping(ctx).Return(errors.New("failed to connect"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewRepository(ctx, &Factory{
				DSN:    cfg.DatabaseDSN,
				Logger: appLogger,
			})

			assert.NoError(t, err)
			assert.NotNil(t, result)
		})
	}
}

func Test_CreateRepository(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	appLogger := logger.NewLogger()

	tests := []struct {
		name     string
		dsn      string
		mockDB   func() (*DatabaseRepository, error)
		expected string
	}{
		{
			name: "Valid DSN, returns DatabaseRepository",
			dsn:  "postgres://postgres:postgres@localhost:5432/shortly-test?sslmode=disable",
			mockDB: func() (*DatabaseRepository, error) {
				return &DatabaseRepository{}, nil
			},
			expected: "*repository.DatabaseRepository",
		},
		{
			name: "Invalid DSN, falls back to InMemoryRepository",
			dsn:  "invalid_dsn",
			mockDB: func() (*DatabaseRepository, error) {
				return nil, errors.New("failed to connect to database")
			},
			expected: "*repository.InMemoryRepository",
		},
		{
			name: "In-Memory repository",
			dsn:  "",
			mockDB: func() (*DatabaseRepository, error) {
				return nil, errors.New("failed to connect")
			},
			expected: "*repository.InMemoryRepository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newDatabaseRepository := tt.mockDB

			originalNewDatabaseRepository := newDatabaseRepository
			defer func() { newDatabaseRepository = originalNewDatabaseRepository }()

			factory := &Factory{
				DSN:    tt.dsn,
				Logger: appLogger,
			}

			repo, err := factory.CreateRepository(ctx)

			assert.NoError(t, err)
			assert.NotNil(t, repo)
			assert.Equal(t, tt.expected, fmt.Sprintf("%T", repo))
		})
	}
}
