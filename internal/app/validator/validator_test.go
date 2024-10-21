package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"shortly/internal/app/errors"
)

func Test_Validate(t *testing.T) {
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
			expected: errors.ErrInvalidURL,
		},
		{
			name:     "Empty URL",
			url:      "",
			valid:    false,
			expected: errors.ErrInvalidURL,
		},
		{
			name:     "URL without scheme",
			url:      "www.example.com",
			valid:    false,
			expected: errors.ErrInvalidURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate(tt.url)

			if tt.valid {
				assert.NoError(t, result)
			} else {
				assert.Error(t, result)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
