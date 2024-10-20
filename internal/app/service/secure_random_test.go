package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
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
				mockSecureRandom.EXPECT().Hex().Return("", errors.New("failed to generate secure random bytes"))
			},
			rand: func() (string, error) {
				return mockSecureRandom.Hex()
			},
			expected: result{
				hex:   "",
				error: errors.New("failed to generate secure random bytes"),
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
