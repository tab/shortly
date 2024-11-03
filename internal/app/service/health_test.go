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
		expected *HealthService
	}{
		{
			name: "Success",
			repo: repo,
			expected: &HealthService{
				repo: repo,
			},
		},
		{
			name: "Nil repository",
			repo: nil,
			expected: &HealthService{
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
	repo := repository.NewMockRepository(ctrl)

	tests := []struct {
		name     string
		before   func()
		expected error
	}{
		{
			name: "Success",
			before: func() {
				repo.EXPECT().Ping(ctx).Return(nil)
			},
			expected: nil,
		},
		{
			name: "Failure",
			before: func() {
				repo.EXPECT().Ping(ctx).Return(errors.New("failed to connect error"))
			},
			expected: errors.New("failed to connect error"),
		},
		{
			name: "Cancelled",
			before: func() {
				repo.EXPECT().Ping(gomock.Any()).Return(context.Canceled)
			},
			expected: context.Canceled,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			service := NewHealthService(repo)
			result := service.Ping(ctx)

			assert.Equal(t, tt.expected, result)
		})
	}
}
