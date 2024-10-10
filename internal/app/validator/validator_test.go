package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"shortly/internal/app/errors"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		valid    bool
		expected error
	}{
		{
			name:     "Valid URL",
			url:      "https://www.google.com",
			valid:    true,
			expected: nil,
		},
		{
			name:     "Invalid URL",
			url:      "not-a-url",
			valid:    false,
			expected: errors.ErrorInvalidURL,
		},
		{
			name:     "Empty URL",
			url:      "",
			valid:    false,
			expected: errors.ErrorInvalidURL,
		},
		{
			name:     "URL without scheme",
			url:      "www.example.com",
			valid:    false,
			expected: errors.ErrorInvalidURL,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Validate(test.url)

			if test.valid {
				assert.NoError(t, result)
			} else {
				assert.Error(t, result)
				assert.Equal(t, test.expected, result)
			}
		})
	}
}
