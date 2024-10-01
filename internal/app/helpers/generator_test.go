package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateShortCode(t *testing.T) {
	type result struct {
		length   int
		fallback bool
	}

	tests := []struct {
		name     string
		expected result
	}{
		{
			name: "Success",
			expected: result{
				length:   8,
				fallback: false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := GenerateShortCode()

			assert.NotEmpty(t, result)
			assert.Equal(t, test.expected.length, len(result))
			assert.NotEqual(t, "randomID", result)
		})
	}
}
