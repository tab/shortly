package worker

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"

	"shortly/internal/app/config"
	"shortly/internal/app/dto"
	"shortly/internal/app/repository"
	"shortly/internal/logger"
)

func Test_worker_StartAndStop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.Background())

	cfg := &config.Config{
		AppEnv: "test",
	}
	repo := repository.NewMockRepository(ctrl)
	appLogger := logger.NewLogger()
	deleteWorker := NewDeleteWorker(ctx, cfg, repo, appLogger)

	assert.NotPanics(t, func() {
		deleteWorker.Start()
	})

	assert.NotPanics(t, func() {
		cancel()
	})
}

func Test_DeleteWorker_Perform(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv: "test",
	}
	repo := repository.NewMockRepository(ctrl)
	appLogger := logger.NewLogger()

	UserUUID, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174001")

	tests := []struct {
		name   string
		params dto.BatchDeleteParams
		before func()
	}{
		{
			name: "Success",
			params: dto.BatchDeleteParams{
				UserID:     UserUUID,
				ShortCodes: []string{"abcd1234"},
			},
			before: func() {
				repo.EXPECT().DeleteURLsByUserID(gomock.Any(), UserUUID, []string{"abcd1234"}).Return(nil)
			},
		},
		{
			name: "Error",
			params: dto.BatchDeleteParams{
				UserID:     UserUUID,
				ShortCodes: []string{"abcd1234"},
			},
			before: func() {
				repo.EXPECT().DeleteURLsByUserID(gomock.Any(), UserUUID, []string{"abcd1234"}).Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			tt.before()

			ctx, cancel := context.WithCancel(context.Background())

			w := NewDeleteWorker(ctx, cfg, repo, appLogger)
			w.Start()
			w.Add(tt.params)

			time.Sleep(50 * time.Millisecond)

			assert.NotPanics(t, func() {
				cancel()
			})
		})
	}
}

func Test_DeleteWorker_Unique(t *testing.T) {
	tests := []struct {
		name       string
		shortCodes []string
		expected   []string
	}{
		{
			name:       "No duplicates",
			shortCodes: []string{"code1", "code2", "code3"},
			expected:   []string{"code1", "code2", "code3"},
		},
		{
			name:       "With duplicates",
			shortCodes: []string{"code1", "code2", "code1", "code3", "code2"},
			expected:   []string{"code1", "code2", "code3"},
		},
		{
			name:       "Empty list",
			shortCodes: []string{},
			expected:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := unique(tt.shortCodes)
			assert.ElementsMatch(t, tt.expected, result)
		})
	}
}
