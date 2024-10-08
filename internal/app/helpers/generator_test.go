package helpers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockSecureRandom struct{}

func (mock *MockSecureRandom) Hex() (string, error) {
	return "abcd1234", nil
}

type MockFailingSecureRandom struct{}

func (mock *MockFailingSecureRandom) Hex() (string, error) {
	return "", errors.New("failed to generate secure random bytes")
}

func TestNewSecureRandomHex(t *testing.T) {
	type result struct {
		code   string
		length int
		error  bool
	}

	testCases := []struct {
		name      string
		generator SecureRandomGenerator
		expected  result
	}{
		{
			name:      "Success",
			generator: NewSecureRandom(),
			expected:  result{code: "unique", length: 8, error: false},
		},
		{
			name:      "Mocked Success",
			generator: &MockSecureRandom{},
			expected:  result{code: "abcd1234", length: 8, error: false},
		},
		{
			name:      "Mocked Failure",
			generator: &MockFailingSecureRandom{},
			expected:  result{code: "", length: 0, error: true},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			code, err := test.generator.Hex()

			if test.expected.error {
				assert.Error(t, err)
				assert.Equal(t, "failed to generate secure random bytes", err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.expected.length, len(code))
		})
	}
}
