package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_InMemoryRepository_Set(t *testing.T) {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.Set(tt.url)
			assert.NoError(t, err)

			storedURL, found := store.Get(tt.url.ShortCode)
			if tt.expected {
				assert.True(t, found)
				assert.Equal(t, tt.url.LongURL, storedURL.LongURL)
			} else {
				assert.False(t, found)
			}
		})
	}
}

func Test_InMemoryRepository_Get(t *testing.T) {
	store := NewInMemoryRepository()

	err := store.Set(URL{
		LongURL:   "https://example.com",
		ShortCode: "abcd1234",
	})
	assert.NoError(t, err)

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			longURL, found := store.Get(tt.shortURL)

			if tt.found {
				assert.NotNil(t, longURL)
				assert.Equal(t, tt.expected, longURL.LongURL)
			} else {
				assert.Nil(t, longURL)
			}
			assert.Equal(t, tt.found, found)
		})
	}
}

func Test_InMemoryRepository_GetAll(t *testing.T) {
	tests := []struct {
		name     string
		before   func(store *InMemoryRepository)
		expected []URL
	}{
		{
			name: "Success",
			before: func(store *InMemoryRepository) {
				err := store.Set(URL{
					UUID:      "6455bd07-e431-4851-af3c-4f703f726639",
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				})
				assert.NoError(t, err)
			},
			expected: []URL{
				{
					UUID:      "6455bd07-e431-4851-af3c-4f703f726639",
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				},
			},
		},
		{
			name: "Multiple URLs",
			before: func(store *InMemoryRepository) {
				err := store.Set(URL{
					UUID:      "6455bd07-e431-4851-af3c-4f703f726639",
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				})
				assert.NoError(t, err)

				err = store.Set(URL{
					UUID:      "3dc48b80-5072-4e23-963c-f5b942ed1a31",
					LongURL:   "https://github.com",
					ShortCode: "ab12ab12",
				})
				assert.NoError(t, err)
			},
			expected: []URL{
				{
					UUID:      "6455bd07-e431-4851-af3c-4f703f726639",
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				},
				{
					UUID:      "3dc48b80-5072-4e23-963c-f5b942ed1a31",
					LongURL:   "https://github.com",
					ShortCode: "ab12ab12",
				},
			},
		},
		{
			name:     "Empty",
			before:   func(_ *InMemoryRepository) {},
			expected: []URL{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewInMemoryRepository()
			tt.before(store)

			results := store.GetAll()

			assert.NotNil(t, results)
			assert.ElementsMatch(t, tt.expected, results)
		})
	}
}
