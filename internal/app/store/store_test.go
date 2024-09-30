package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	store := NewURLStore()
	store.Set("GitHub", "HTTPS://GITHUB.COM")

	tests := []struct {
		name     string
		shortURL string
		longURL  string
	}{
		{
			name:     "Success",
			shortURL: "abcd1234",
			longURL:  "https://example.com",
		},
		{
			name:     "Overwrite",
			shortURL: "GitHub",
			longURL:  "https://github.com",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store.Set(test.shortURL, test.longURL)

			assert.NoError(t, nil)
		})
	}
}

func TestGet(t *testing.T) {
	store := NewURLStore()
	store.Set("abcd1234", "https://example.com")

	tests := []struct {
		name     string
		shortURL string
		expected string
		found    bool
	}{
		{
			name:     "Success",
			shortURL: "abcd1234",
			expected: "https://example.com",
			found:    true,
		},
		{
			name:     "Not Found",
			shortURL: "1234abcd",
			expected: "",
			found:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			longURL, found := store.Get(test.shortURL)

			assert.Equal(t, test.expected, longURL)
			assert.Equal(t, test.found, found)
		})
	}
}
