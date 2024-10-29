package persistence

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortly/internal/app/errors"
	"shortly/internal/app/repository"
	"shortly/internal/logger"
)

func Test_PersistenceManager_Load(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepository(ctrl)
	mockFileRepo := repository.NewMockFileRepository(ctrl)
	appLogger := logger.NewLogger()

	snapshot := &repository.Memento{
		State: []repository.URL{
			{
				UUID:      "6455bd07-e431-4851-af3c-4f703f726639",
				LongURL:   "http://example.com",
				ShortCode: "abcd1234",
			},
		},
	}

	tests := []struct {
		name   string
		before func()
	}{
		{
			name: "Success",
			before: func() {
				mockFileRepo.EXPECT().Load().Return(snapshot, nil)
				mockRepo.EXPECT().Restore(snapshot)
			},
		},
		{
			name: "Failure",
			before: func() {
				mockFileRepo.EXPECT().Load().Return(nil, errors.ErrFailedToOpenFile)
				appLogger.Error().Err(errors.ErrFailedToOpenFile).Msg("Failed to load data from file")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			pm := NewPersistenceManager(mockRepo, mockFileRepo, appLogger)
			err := pm.Load()

			assert.NoError(t, err)
		})
	}
}

func Test_PersistenceManager_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockRepository(ctrl)
	mockFileRepo := repository.NewMockFileRepository(ctrl)
	appLogger := logger.NewLogger()

	snapshot := &repository.Memento{
		State: []repository.URL{
			{
				UUID:      "6455bd07-e431-4851-af3c-4f703f726639",
				LongURL:   "http://example.com",
				ShortCode: "abcd1234",
			},
		},
	}

	tests := []struct {
		name   string
		before func()
	}{
		{
			name: "Success",
			before: func() {
				mockRepo.EXPECT().CreateMemento().Return(snapshot)
				mockFileRepo.EXPECT().Save(snapshot).Return(nil)
			},
		},
		{
			name: "Failure",
			before: func() {
				mockRepo.EXPECT().CreateMemento().Return(snapshot)
				mockFileRepo.EXPECT().Save(snapshot).Return(errors.ErrFailedToOpenFile)
				appLogger.Error().Err(errors.ErrFailedToOpenFile).Msg("Failed to save data to file")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			pm := NewPersistenceManager(mockRepo, mockFileRepo, appLogger)
			err := pm.Save()

			assert.NoError(t, err)
		})
	}
}
