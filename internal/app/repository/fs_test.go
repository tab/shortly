package repository

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FileStorageRepository_Set(t *testing.T) {
	tests := []struct {
		name     string
		url      URL
		expected bool
	}{
		{
			name: "Add new URL",
			url: URL{
				UUID:      "6455bd07-e431-4851-af3c-4f703f726639",
				LongURL:   "https://example.com",
				ShortCode: "abcd1234",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			filePath := t.TempDir() + "/store-test.json"

			store, err := NewFileStorageRepository(ctx, filePath)
			assert.NoError(t, err)

			err = store.Set(tt.url)
			assert.NoError(t, err)

			storedURL, found := store.Get(tt.url.ShortCode)
			if tt.expected {
				assert.True(t, found)
				assert.Equal(t, tt.url.LongURL, storedURL.LongURL)
			} else {
				assert.False(t, found)
			}

			cancel()
			store.wg.Wait()

			t.Cleanup(func() {
				os.RemoveAll(filePath)
			})
		})
	}
}

func Test_FileStorageRepository_Get(t *testing.T) {
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
			ctx, cancel := context.WithCancel(context.Background())

			filePath := t.TempDir() + "/store-test.json"

			store, err := NewFileStorageRepository(ctx, filePath)
			assert.NoError(t, err)

			err = store.Set(URL{
				UUID:      "6455bd07-e431-4851-af3c-4f703f726639",
				LongURL:   "https://example.com",
				ShortCode: "abcd1234",
			})
			assert.NoError(t, err)

			longURL, found := store.Get(tt.shortURL)

			if tt.found {
				assert.NotNil(t, longURL)
				assert.Equal(t, tt.expected, longURL.LongURL)
			} else {
				assert.Nil(t, longURL)
			}
			assert.Equal(t, tt.found, found)

			cancel()
			store.wg.Wait()

			t.Cleanup(func() {
				os.RemoveAll(filePath)
			})
		})
	}
}

func Test_FileStorageRepository_GetAll(t *testing.T) {
	tests := []struct {
		name     string
		before   func(store *FileStorageRepository)
		expected []URL
	}{
		{
			name: "Success",
			before: func(store *FileStorageRepository) {
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
			before: func(store *FileStorageRepository) {
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
			before:   func(_ *FileStorageRepository) {},
			expected: []URL{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			filePath := t.TempDir() + "/store-test.json"

			store, err := NewFileStorageRepository(ctx, filePath)
			assert.NoError(t, err)

			tt.before(store)

			results := store.GetAll()

			assert.NotNil(t, results)
			assert.ElementsMatch(t, tt.expected, results)

			cancel()
			store.wg.Wait()

			t.Cleanup(func() {
				os.RemoveAll(filePath)
			})
		})
	}
}

func Test_FileStorageRepository_Load(t *testing.T) {
	tests := []struct {
		name     string
		before   func(store *FileStorageRepository)
		expected error
	}{
		{
			name: "Success",
			before: func(store *FileStorageRepository) {
				err := store.Set(URL{
					UUID:      "6455bd07-e431-4851-af3c-4f703f726639",
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				})
				assert.NoError(t, err)

				err = store.Save()
				assert.NoError(t, err)
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			filePath := t.TempDir() + "/store-test.json"

			store, err := NewFileStorageRepository(ctx, filePath)
			assert.NoError(t, err)

			tt.before(store)

			result := store.Load()
			assert.Equal(t, tt.expected, result)

			cancel()
			store.wg.Wait()

			t.Cleanup(func() {
				os.RemoveAll(filePath)
			})
		})
	}
}

func Test_FileStorageRepository_Save(t *testing.T) {
	tests := []struct {
		name     string
		before   func(store *FileStorageRepository)
		expected error
	}{
		{
			name: "Success",
			before: func(store *FileStorageRepository) {
				err := store.Set(URL{
					UUID:      "6455bd07-e431-4851-af3c-4f703f726639",
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				})
				assert.NoError(t, err)
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			filePath := t.TempDir() + "/store-test.json"

			store, err := NewFileStorageRepository(ctx, filePath)
			assert.NoError(t, err)

			tt.before(store)

			result := store.Save()
			assert.Equal(t, tt.expected, result)

			cancel()
			store.wg.Wait()

			t.Cleanup(func() {
				os.RemoveAll(filePath)
			})
		})
	}
}

func Test_FileStorageRepository_Wait(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(store *FileStorageRepository)
		expected []URL
	}{
		{
			name: "Wait for snapshot to save on context cancellation",
			setup: func(store *FileStorageRepository) {
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			filePath := t.TempDir() + "/store-test.json"

			store, err := NewFileStorageRepository(ctx, filePath)
			tt.setup(store)
			assert.NoError(t, err)

			cancel()
			store.Wait()

			results := store.GetAll()
			assert.ElementsMatch(t, tt.expected, results)

			t.Cleanup(func() {
				os.RemoveAll(filePath)
			})
		})
	}
}
