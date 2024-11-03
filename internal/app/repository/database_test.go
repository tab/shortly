package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func Test_DatabaseRepository_Ping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPool := NewMockPgxPool(ctrl)
	mockRow := NewMockRow(ctrl)

	dbRepo := &DatabaseRepository{db: mockPool}
	ctx := context.Background()

	tests := []struct {
		name     string
		before   func()
		expected error
	}{
		{
			name: "Success",
			before: func() {
				mockPool.EXPECT().QueryRow(ctx, "SELECT 1").Return(mockRow)
				mockRow.EXPECT().Scan(gomock.Any()).Return(nil)
			},
			expected: nil,
		},
		{
			name: "Failed",
			before: func() {
				mockPool.EXPECT().QueryRow(ctx, "SELECT 1").Return(mockRow)
				mockRow.EXPECT().Scan(gomock.Any()).Return(errors.New("failed to connect"))
			},
			expected: errors.New("failed to connect"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			err := dbRepo.Ping(ctx)

			if tt.expected != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_DatabaseRepository_Set(t *testing.T) {
	dbRepo := &DatabaseRepository{}

	err := dbRepo.Set(URL{})
	assert.NoError(t, err)
}

func Test_DatabaseRepository_Get(t *testing.T) {
	dbRepo := &DatabaseRepository{}

	url, found := dbRepo.Get("abcd1234")
	assert.False(t, found)
	assert.Nil(t, url)
}

func Test_DatabaseRepository_CreateMemento(t *testing.T) {
	dbRepo := &DatabaseRepository{}

	memento := dbRepo.CreateMemento()
	assert.NotNil(t, memento)
	assert.Empty(t, memento.State)
}
