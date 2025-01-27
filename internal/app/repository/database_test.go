package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shortly/internal/app/repository/db"
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

func Benchmark_DatabaseRepository_CreateURL(b *testing.B) {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_DSN")
	store, err := NewDatabaseRepository(ctx, dsn)
	assert.NoError(b, err)

	b.Cleanup(func() {
		err = spec.TruncateTables(ctx, dsn)
		require.NoError(b, err)

		store.Close()
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		UUID, err := uuid.NewRandom()
		require.NoError(b, err)

		longURL := fmt.Sprintf("https://example.com/%d", i)
		shortCode := fmt.Sprintf("abcd%d", i)

		_, err = store.CreateURL(ctx, URL{
			UUID:      UUID,
			LongURL:   longURL,
			ShortCode: shortCode,
		})
		assert.NoError(b, err)
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

func Benchmark_DatabaseRepository_GetURLByShortCode(b *testing.B) {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_DSN")
	store, err := NewDatabaseRepository(ctx, dsn)
	assert.NoError(b, err)

	UUID, err := uuid.NewRandom()
	require.NoError(b, err)

	longURL := "https://example.com"
	shortCode := "abcd1234"

	_, err = store.CreateURL(ctx, URL{
		UUID:      UUID,
		LongURL:   longURL,
		ShortCode: shortCode,
	})
	assert.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, found := store.GetURLByShortCode(ctx, shortCode)
		assert.True(b, found)
	}
}

func Test_DatabaseRepository_GetURLsByUserID(t *testing.T) {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_DSN")
	store, err := NewDatabaseRepository(ctx, dsn)
	assert.NoError(t, err)

	UUID1, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f720001")
	UUID2, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f720002")
	UserUUID1, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174001")
	UserUUID2, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174002")

	type result struct {
		count     int
		UUID      uuid.UUID
		LongURL   string
		ShortCode string
	}

	tests := []struct {
		name     string
		before   func()
		UserID   uuid.UUID
		expected result
	}{
		{
			name: "Success",
			before: func() {
				_, err = store.CreateURL(ctx, URL{
					UUID:      UUID1,
					LongURL:   "https://google.com",
					ShortCode: "abcd0001",
				})
				assert.NoError(t, err)

				_, err = store.CreateURL(ctx, URL{
					UUID:      UUID2,
					LongURL:   "https://github.com",
					ShortCode: "abcd0002",
					UserUUID:  UserUUID1,
				})
				assert.NoError(t, err)
			},
			UserID: UserUUID1,
			expected: result{
				count:     1,
				UUID:      UUID2,
				LongURL:   "https://github.com",
				ShortCode: "abcd0002",
			},
		},
		{
			name: "Not owned",
			before: func() {
				_, err = store.CreateURL(ctx, URL{
					UUID:      UUID1,
					LongURL:   "https://google.com",
					ShortCode: "abcd0001",
					UserUUID:  UserUUID2,
				})
				assert.NoError(t, err)
			},
			UserID: UserUUID1,
			expected: result{
				count:     0,
				UUID:      uuid.Nil,
				LongURL:   "",
				ShortCode: "",
			},
		},
		{
			name: "Not Found",
			before: func() {
				_, err = store.CreateURL(ctx, URL{
					UUID:      UUID1,
					LongURL:   "https://google.com",
					ShortCode: "abcd0001",
				})
				assert.NoError(t, err)
			},
			UserID: UserUUID1,
			expected: result{
				count:     0,
				UUID:      uuid.Nil,
				LongURL:   "",
				ShortCode: "",
			},
		},
		{
			name:   "Empty",
			before: func() {},
			UserID: uuid.Nil,
			expected: result{
				count:     0,
				UUID:      uuid.Nil,
				LongURL:   "",
				ShortCode: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			rows, total, err := store.GetURLsByUserID(ctx, tt.UserID, 25, 0)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.count, total)

			if tt.expected.count > 0 {
				assert.NotEmpty(t, rows)
				assert.Equal(t, tt.expected.UUID, rows[0].UUID)
				assert.Equal(t, tt.expected.LongURL, rows[0].LongURL)
				assert.Equal(t, tt.expected.ShortCode, rows[0].ShortCode)
			} else {
				assert.Empty(t, rows)
			}

			t.Cleanup(func() {
				err = spec.TruncateTables(ctx, dsn)
				require.NoError(t, err)
			})
		})
	}
}

func Test_DatabaseRepository_DeleteURLsByUserID(t *testing.T) {
	ctx := context.Background()
	dsn := os.Getenv("DATABASE_DSN")
	store, err := NewDatabaseRepository(ctx, dsn)
	assert.NoError(t, err)

	UUID1, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f720001")
	UUID2, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f720002")
	UserUUID1, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174001")
	UserUUID2, _ := uuid.Parse("123e4567-e89b-12d3-a456-426614174002")

	tests := []struct {
		name          string
		currentUserID uuid.UUID
		ownerID       uuid.UUID
		before        func()
		params        db.DeleteURLsByUserIDAndShortCodesParams
		expected      int
	}{
		{
			name:          "Success",
			currentUserID: UserUUID1,
			ownerID:       UserUUID1,
			before: func() {
				_, err = store.CreateURL(ctx, URL{
					UUID:      UUID1,
					LongURL:   "https://google.com",
					ShortCode: "abcd0001",
					UserUUID:  UserUUID1,
				})
				assert.NoError(t, err)

				_, err = store.CreateURL(ctx, URL{
					UUID:      UUID2,
					LongURL:   "https://github.com",
					ShortCode: "abcd0002",
					UserUUID:  UserUUID1,
				})
				assert.NoError(t, err)
			},
			params: db.DeleteURLsByUserIDAndShortCodesParams{
				UserUUID:   UserUUID1,
				ShortCodes: []string{"abcd0001", "abcd0002"},
			},
			expected: 0,
		},
		{
			name:          "Not owned",
			currentUserID: UserUUID1,
			ownerID:       UserUUID2,
			before: func() {
				_, err = store.CreateURL(ctx, URL{
					UUID:      UUID1,
					LongURL:   "https://google.com",
					ShortCode: "abcd0001",
					UserUUID:  UserUUID2,
				})
				assert.NoError(t, err)
			},
			params: db.DeleteURLsByUserIDAndShortCodesParams{
				UserUUID:   UserUUID1,
				ShortCodes: []string{"abcd0001"},
			},
			expected: 1,
		},
		{
			name:          "Not found",
			currentUserID: UserUUID1,
			ownerID:       UserUUID1,
			before: func() {
				_, err = store.CreateURL(ctx, URL{
					UUID:      UUID1,
					LongURL:   "https://google.com",
					ShortCode: "abcd0001",
					UserUUID:  UserUUID1,
				})
				assert.NoError(t, err)
			},
			params: db.DeleteURLsByUserIDAndShortCodesParams{
				UserUUID:   UserUUID1,
				ShortCodes: []string{"1234abcd"},
			},
			expected: 1,
		},
		{
			name:          "Empty",
			currentUserID: UserUUID1,
			ownerID:       UserUUID1,
			before: func() {
				_, err = store.CreateURL(ctx, URL{
					UUID:      UUID1,
					LongURL:   "https://google.com",
					ShortCode: "abcd0001",
					UserUUID:  UserUUID1,
				})
				assert.NoError(t, err)
			},
			params: db.DeleteURLsByUserIDAndShortCodesParams{
				UserUUID:   UserUUID1,
				ShortCodes: []string{},
			},
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			err = store.DeleteURLsByUserID(ctx, tt.params.UserUUID, tt.params.ShortCodes)
			assert.NoError(t, err)

			_, total, err := store.GetURLsByUserID(ctx, tt.ownerID, 25, 0)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, total)

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
