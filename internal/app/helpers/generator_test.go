package helpers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockSecureRandom struct{}

func (MockSecureRandom) Hex() (string, error) {
	return "abcd1234", nil
}

type MockFailingGenerator struct{}

func (MockFailingGenerator) Hex() (string, error) {
	return "", errors.New("failed to generate short code")
}

func TestSecureRandom_Hex(t *testing.T) {
	type result struct {
		length int
		code   string
		error  bool
	}

	testCases := []struct {
		name     string
		mock     SecureRandomGenerator
		expected result
	}{
		{
			name: "Success",
			mock: NewSecureRandom(),
			expected: result{
				length: BytesLength * 2,
				code:   "",
				error:  false,
			},
		},
		{
			name: "Mocked Success",
			mock: MockSecureRandom{},
			expected: result{
				length: BytesLength * 2,
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
			code, err := test.mock.Hex()

			if test.expected.error {
				assert.Error(t, err)
				assert.Empty(t, code)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, code)
				assert.Equal(t, test.expected.length, len(code))

				if test.expected.code != "" {
					assert.Equal(t, test.expected.code, code)
				}
			}
		})
	}
}
