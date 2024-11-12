package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortly/internal/app/repository"
)

func Test_NewHealthService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repository.NewMockRepository(ctrl)

	tests := []struct {
		name     string
		repo     repository.Repository
		expected *healthService
	}{
		{
			name: "Success",
			repo: repo,
			expected: &healthService{
				repo: repo,
			},
		},
		{
			name: "Nil repository",
			repo: nil,
			expected: &healthService{
				repo: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewHealthService(tt.repo)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_HealthService_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockDatabaseRepo := repository.NewMockDatabase(ctrl)
	inMemoryRepo := repository.NewInMemoryRepository()

	tests := []struct {
		name     string
		repo     repository.Repository
		before   func()
		expected error
	}{
		{
			name: "Success",
			repo: mockDatabaseRepo,
			before: func() {
				mockDatabaseRepo.EXPECT().Ping(ctx).Return(nil)
			},
			expected: nil,
		},
		{
			name: "Failure",
			repo: mockDatabaseRepo,
			before: func() {
				mockDatabaseRepo.EXPECT().Ping(ctx).Return(errors.New("failed to connect error"))
			},
			expected: errors.New("failed to connect error"),
		},
		{
			name: "Cancelled",
			repo: mockDatabaseRepo,
			before: func() {
				mockDatabaseRepo.EXPECT().Ping(gomock.Any()).Return(context.Canceled)
			},
			expected: context.Canceled,
		},
		{
			name:     "InMemoryRepository",
			repo:     inMemoryRepo,
			before:   func() {},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			service := NewHealthService(tt.repo)
			result := service.Ping(ctx)

			assert.Equal(t, tt.expected, result)
		})
	}
}
