package repository

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shortly/internal/spec"
)

func TestMain(m *testing.M) {
	if err := spec.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	if os.Getenv("GO_ENV") == "ci" {
		os.Exit(0)
	}

	code := m.Run()
	os.Exit(code)
}

func Test_DatabaseRepository_CreateURL(t *testing.T) {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_DSN")
	store, err := NewDatabaseRepository(ctx, dsn)
	assert.NoError(t, err)

	UUID, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f726639")

	tests := []struct {
		name       string
		before     func()
		attributes URL
		expected   URL
	}{
		{
			name: "Success",
			attributes: URL{
				UUID:      UUID,
				LongURL:   "https://example.com",
				ShortCode: "abcd1234",
			},
			expected: URL{
				UUID:      UUID,
				LongURL:   "https://example.com",
				ShortCode: "abcd1234",
			},
		},
		{
			name: "Not unique",
			before: func() {
				_, err := store.CreateURL(ctx, URL{
					UUID:      UUID,
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
				})
				assert.NoError(t, err)
			},
			attributes: URL{
				UUID:      UUID,
				LongURL:   "https://example.com",
				ShortCode: "abcd0001",
			},
			expected: URL{
				UUID:      UUID,
				LongURL:   "https://example.com",
				ShortCode: "abcd1234",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row, err := store.CreateURL(ctx, tt.attributes)
			assert.NoError(t, err)
			assert.Equal(t, tt.attributes.LongURL, row.LongURL)

			storedURL, found := store.GetURLByShortCode(ctx, tt.attributes.ShortCode)
			assert.True(t, found)
			assert.Equal(t, tt.attributes.UUID, storedURL.UUID)
			assert.Equal(t, tt.attributes.LongURL, storedURL.LongURL)
			assert.Equal(t, tt.attributes.ShortCode, storedURL.ShortCode)

			t.Cleanup(func() {
				err = spec.TruncateTables(ctx, dsn)
				require.NoError(t, err)
			})
		})
	}
}

func Test_DatabaseRepository_CreateURLs(t *testing.T) {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_DSN")
	store, err := NewDatabaseRepository(ctx, dsn)
	assert.NoError(t, err)

	UUID1, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f720001")
	UUID2, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f720002")

	tests := []struct {
		name     string
		urls     []URL
		expected bool
	}{
		{
			name: "Success",
			urls: []URL{
				{
					UUID:      UUID1,
					LongURL:   "https://google.com",
					ShortCode: "abcd0001",
				},
				{
					UUID:      UUID2,
					LongURL:   "https://github.com",
					ShortCode: "abcd0002",
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = store.CreateURLs(ctx, tt.urls)
			assert.NoError(t, err)

			for _, url := range tt.urls {
				storedURL, found := store.GetURLByShortCode(ctx, url.ShortCode)
				if tt.expected {
					assert.True(t, found)
					assert.Equal(t, url.LongURL, storedURL.LongURL)
				} else {
					assert.False(t, found)
				}
			}

			t.Cleanup(func() {
				err = spec.TruncateTables(ctx, dsn)
				require.NoError(t, err)
			})
		})
	}
}

func Test_DatabaseRepository_GetURLByShortCode(t *testing.T) {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_DSN")
	store, err := NewDatabaseRepository(ctx, dsn)
	assert.NoError(t, err)

	UUID, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f726639")

	_, err = store.CreateURL(ctx, URL{
		UUID:      UUID,
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
			row, found := store.GetURLByShortCode(ctx, tt.shortURL)

			if tt.found {
				assert.NotNil(t, row)
				assert.Equal(t, tt.expected, row.LongURL)
			} else {
				assert.Nil(t, row)
			}
			assert.Equal(t, tt.found, found)

			t.Cleanup(func() {
				err = spec.TruncateTables(ctx, dsn)
				require.NoError(t, err)
			})
		})
	}
}

func Test_DatabaseRepository_Ping(t *testing.T) {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_DSN")
	store, err := NewDatabaseRepository(ctx, dsn)
	assert.NoError(t, err)

	err = store.Ping(ctx)
	assert.NoError(t, err)
}

func TestDatabaseRepo_Close(t *testing.T) {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_DSN")
	store, err := NewDatabaseRepository(ctx, dsn)
	assert.NoError(t, err)

	store.Close()
}