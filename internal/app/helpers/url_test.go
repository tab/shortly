package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "Valid URL",
			url:      "https://www.google.com",
			expected: true,
		},
		{
			name:     "Invalid URL",
			url:      "not-a-url",
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsValidURL(test.url)

			assert.Equal(t, test.expected, result)
		})
	}
}

func TestIsInvalidURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "Valid URL",
			url:      "https://www.google.com",
			expected: false,
		},
		{
			name:     "Invalid URL",
			url:      "not-a-url",
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := IsInvalidURL(test.url)

			assert.Equal(t, test.expected, result)
		})
	}
}
