package repository

import (
	"context"
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
	cfg := config.LoadConfig()
	appLogger := logger.NewLogger()

	tests := []struct {
		name         string
		dsn          string
		repo         Repository
		expectedType interface{}
	}{
		{
			name:         "In-Memory repository",
			dsn:          "",
			expectedType: &InMemoryRepo{},
		},
		{
			name:         "PostgreSQL database repository",
			dsn:          cfg.DatabaseDSN,
			expectedType: &DatabaseRepo{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewRepository(ctx, &Factory{
				DSN:    tt.dsn,
				Logger: appLogger,
			})

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.IsType(t, tt.expectedType, result)
		})
	}
}
