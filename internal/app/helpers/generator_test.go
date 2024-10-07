package helpers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockShortCode struct{}

func (MockShortCode) Code() (string, error) {
	return "abcd1234", nil
}

type MockFailingGenerator struct{}

func (MockFailingGenerator) Code() (string, error) {
	return "", errors.New("failed to generate short code")
}

func TestShortCode(t *testing.T) {
	type result struct {
		length int
		code   string
		error  bool
	}

	testCases := []struct {
		name     string
		mock     SecureRandom
		expected result
	}{
		{
			name: "Success",
			mock: Generator{},
			expected: result{
				length: ShortCodeLength,
				code:   "unique",
				error:  false,
			},
		},
		{
			name: "Mocked Success",
			mock: MockShortCode{},
			expected: result{
				length: ShortCodeLength,
				code:   "abcd1234",
				error:  false,
			},
		},
		{
			name: "Mocked Failure",
			mock: MockFailingGenerator{},
			expected: result{
				length: 0,
				code:   "",
				error:  true,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			SetShortCodeGenerator(test.mock)

			code, err := ShortCode()

			if test.expected.error {
				assert.Error(t, err)
				assert.Empty(t, code)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, code)
				assert.Equal(t, test.expected.length, len(code))
			}
		})
	}
}
