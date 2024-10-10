package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	store := NewInMemoryRepository()

	tests := []struct {
		name     string
		url      URL
		before   func()
		expected bool
	}{
		{
			name: "Add new URL",
			url: URL{
				LongURL:   "https://example.com",
				ShortCode: "abcd1234",
			},
			before:   func() {},
			expected: true,
		},
		{
			name: "Overwrite existing URL",
			url: URL{
				LongURL:   "https://github.com",
				ShortCode: "GitHub",
			},
			before: func() {
				store.Set(URL{
					LongURL:   "https://example.com",
					ShortCode: "123456ab",
				})
			},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			store.Set(test.url)

			storedURL, found := store.Get(test.url.ShortCode)
			if test.expected {
				assert.True(t, found)
				assert.Equal(t, test.url.LongURL, storedURL.LongURL)
			} else {
				assert.False(t, found)
			}
		})
	}
}

func TestGet(t *testing.T) {
	store := NewInMemoryRepository()
	store.Set(URL{
		LongURL:   "https://example.com",
		ShortCode: "abcd1234",
	})

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
		{
			name:     "Empty",
			shortURL: "",
			expected: "",
			found:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			longURL, found := store.Get(test.shortURL)

			if test.found {
				assert.NotNil(t, longURL)
				assert.Equal(t, test.expected, longURL.LongURL)
			} else {
				assert.Nil(t, longURL)
			}
			assert.Equal(t, test.found, found)
		})
	}
}
