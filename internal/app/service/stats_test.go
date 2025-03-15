package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortly/internal/app/repository"
)

func Test_NewStatsReporter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := repository.NewMockRepository(ctrl)

	tests := []struct {
		name     string
		repo     repository.Repository
		expected *statsService
	}{
		{
			name: "Success",
			repo: repo,
			expected: &statsService{
				repo: repo,
			},
		},
		{
			name: "Nil repository",
			repo: nil,
			expected: &statsService{
				repo: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewStatsReporter(tt.repo)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_StatsReporter_Counters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	repo := repository.NewMockDatabase(ctrl)

	type result struct {
		urls  int
		users int
	}

	tests := []struct {
		name     string
		before   func()
		repo     repository.Repository
		expected result
		err      error
	}{
		{
			name: "Success",
			before: func() {
				repo.EXPECT().Counters(ctx).Return(10, 2, nil)
			},
			repo: repo,
			expected: result{
				urls:  10,
				users: 2,
			},
		},
		{
			name: "Empty",
			before: func() {
				repo.EXPECT().Counters(ctx).Return(0, 0, nil)
			},
			repo: repo,
			expected: result{
				urls:  0,
				users: 0,
			},
		},
		{
			name: "Error",
			before: func() {
				repo.EXPECT().Counters(ctx).Return(0, 0, assert.AnError)
			},
			repo: repo,
			expected: result{
				urls:  0,
				users: 0,
			},
			err: assert.AnError,
		},
		{
			name:   "InMemory repository",
			before: func() {},
			repo:   repository.NewMockInMemory(ctrl),
			expected: result{
				urls:  0,
				users: 0,
			},
		},
		{
			name:   "Nil repository",
			before: func() {},
			repo:   nil,
			expected: result{
				urls:  0,
				users: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			service := NewStatsReporter(tt.repo)
			urls, users, err := service.Counters(ctx)

			if tt.err != nil {
				assert.Equal(t, tt.err, err)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, tt.expected.urls, urls)
			assert.Equal(t, tt.expected.users, users)
		})
	}
}
