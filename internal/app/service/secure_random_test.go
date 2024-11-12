package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"shortly/internal/app/errors"
)

func Test_SecureRandom_Hex(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSecureRandom := NewMockSecureRandomGenerator(ctrl)
	secureRandom := NewSecureRandom()

	type result struct {
		hex   string
		error error
	}

	tests := []struct {
		name     string
		mocked   bool
		before   func()
		rand     func() (string, error)
		expected result
	}{
		{
			name:   "Success",
			mocked: false,
			before: func() {},
			rand: func() (string, error) {
				return secureRandom.Hex()
			},
			expected: result{
				hex:   "random-hex-string",
				error: nil,
			},
		},
		{
			name:   "Mocked success",
			mocked: true,
			before: func() {
				mockSecureRandom.EXPECT().Hex().Return("abcd1234", nil)
			},
			rand: func() (string, error) {
				return mockSecureRandom.Hex()
			},
			expected: result{
				hex:   "abcd1234",
				error: nil,
			},
		},
		{
			name:   "Mocked failure",
			mocked: true,
			before: func() {
				mockSecureRandom.EXPECT().Hex().Return("", errors.ErrFailedToReadRandomBytes)
			},
			rand: func() (string, error) {
				return mockSecureRandom.Hex()
			},
			expected: result{
				hex:   "",
				error: errors.ErrFailedToReadRandomBytes,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			hex, err := tt.rand()

			if tt.expected.error != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expected.error, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.mocked {
				assert.Equal(t, tt.expected.hex, hex)
			} else {
				assert.NotEmpty(t, hex)
			}
		})
	}
}

func Test_SecureRandom_UUID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSecureRandom := NewMockSecureRandomGenerator(ctrl)
	secureRandom := NewSecureRandom()

	UUID, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f726639")
	generatedUUID, _ := uuid.NewRandom()

	type result struct {
		uuid  uuid.UUID
		error error
	}

	tests := []struct {
		name     string
		mocked   bool
		before   func()
		rand     func() (uuid.UUID, error)
		expected result
	}{
		{
			name:   "Success",
			mocked: false,
			before: func() {},
			rand: func() (uuid.UUID, error) {
				return secureRandom.UUID()
			},
			expected: result{
				uuid:  generatedUUID,
				error: nil,
			},
		},
		{
			name:   "Mocked success",
			mocked: true,
			before: func() {
				mockSecureRandom.EXPECT().UUID().Return(UUID, nil)
			},
			rand: func() (uuid.UUID, error) {
				return mockSecureRandom.UUID()
			},
			expected: result{
				uuid:  UUID,
				error: nil,
			},
		},
		{
			name:   "Mocked failure",
			mocked: true,
			before: func() {
				mockSecureRandom.EXPECT().UUID().Return(uuid.UUID{}, errors.ErrFailedToGenerateUUID)
			},
			rand: func() (uuid.UUID, error) {
				return mockSecureRandom.UUID()
			},
			expected: result{
				uuid:  uuid.UUID{},
				error: errors.ErrFailedToGenerateUUID,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			result, err := tt.rand()

			if tt.expected.error != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expected.error, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.mocked {
				assert.Equal(t, tt.expected.uuid, result)
			} else {
				assert.NotEmpty(t, result)
			}
		})
	}
}
