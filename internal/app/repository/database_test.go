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

func Test_DatabaseRepository_Set(t *testing.T) {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_DSN")
	store, err := NewDatabaseRepository(ctx, dsn)
	assert.NoError(t, err)

	UUID, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f726639")

	tests := []struct {
		name     string
		url      URL
		expected bool
	}{
		{
			name: "Success",
			url: URL{
				UUID:      UUID,
				LongURL:   "https://example.com",
				ShortCode: "abcd1234",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = store.Set(ctx, tt.url)
			assert.NoError(t, err)

			storedURL, found := store.Get(ctx, tt.url.ShortCode)
			if tt.expected {
				assert.True(t, found)
				assert.Equal(t, tt.url.LongURL, storedURL.LongURL)
			} else {
				assert.False(t, found)
			}

			t.Cleanup(func() {
				err = spec.TruncateTables(ctx, dsn)
				require.NoError(t, err)
			})
		})
	}
}

func Test_DatabaseRepository_Get(t *testing.T) {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_DSN")
	store, err := NewDatabaseRepository(ctx, dsn)
	assert.NoError(t, err)

	UUID, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f726639")

	err = store.Set(ctx, URL{
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
			longURL, found := store.Get(ctx, tt.shortURL)

			if tt.found {
				assert.NotNil(t, longURL)
				assert.Equal(t, tt.expected, longURL.LongURL)
			} else {
				assert.Nil(t, longURL)
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
