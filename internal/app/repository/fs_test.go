package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"shortly/internal/app/errors"
)

func Test_FileStorageRepository_Set(t *testing.T) {
	filePath := os.TempDir() + "/store.json"
	store := NewFileStorageRepository(filePath)

	tests := []struct {
		name     string
		url      URL
		before   func()
		expected bool
		error    error
	}{
		{
			name: "Success",
			url: URL{
				UUID:      "6455bd07-e431-4851-af3c-4f703f726639",
				LongURL:   "https://example.com",
				ShortCode: "abcd1234",
			},
			before:   func() {},
			expected: true,
			error:    nil,
		},
		{
			name: "Failed to open file",
			url: URL{
				UUID:      "6455bd07-e431-4851-af3c-4f703f726639",
				LongURL:   "https://example.com",
				ShortCode: "abcd1234",
			},
			before: func() {
				err := os.Mkdir(filePath, 0755)
				assert.NoError(t, err)
			},
			expected: false,
			error:    errors.ErrFailedToOpenFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			err := store.Set(tt.url)

			if tt.error != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			storedURL, found := store.Get(tt.url.ShortCode)
			if tt.expected {
				assert.True(t, found)
				assert.Equal(t, tt.url.UUID, storedURL.UUID)
				assert.Equal(t, tt.url.LongURL, storedURL.LongURL)
			} else {
				assert.False(t, found)
			}

			t.Cleanup(func() {
				os.Remove(filePath)
			})
		})
	}
}

func Test_FileStorageRepository_Get(t *testing.T) {
	filePath := os.TempDir() + "/store.json"
	store := NewFileStorageRepository(filePath)

	err := store.Set(URL{
		UUID:      "6455bd07-e431-4851-af3c-4f703f726639",
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

			t.Cleanup(func() {
				os.Remove(filePath)
			})
		})
	}
}
