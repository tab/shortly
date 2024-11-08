package repository

import (
	"context"
	"errors"
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
