package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"shortly/internal/app/repository/db"
)

func Test_InMemoryRepository_CreateURL(t *testing.T) {
	ctx := context.Background()
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
				store.CreateURL(ctx, URL{
					LongURL:   "https://example.com",
					ShortCode: "123456ab",
				})
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := store.CreateURL(ctx, tt.url)
			assert.NoError(t, err)

			storedURL, found := store.GetURLByShortCode(ctx, tt.url.ShortCode)
			if tt.expected {
				assert.True(t, found)
				assert.Equal(t, tt.url.LongURL, storedURL.LongURL)
			} else {
				assert.False(t, found)
			}
		})
	}
}

func Benchmark_InMemoryRepository_CreateURL(b *testing.B) {
	ctx := context.Background()
	store := NewInMemoryRepository()

	UUID := uuid.MustParse("6455bd07-e431-4851-af3c-4f703f726639")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		longURL := fmt.Sprintf("https://example.com/%d", i)
		shortCode := fmt.Sprintf("abcd%d", i)

		_, err := store.CreateURL(ctx, URL{
			UUID:      UUID,
			LongURL:   longURL,
			ShortCode: shortCode,
		})
		assert.NoError(b, err)
	}
}

func Test_InMemoryRepository_CreateURLs(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryRepository()

	tests := []struct {
		name     string
		urls     []URL
		expected int
	}{
		{
			name: "Add new URLs",
			urls: []URL{
				{
					LongURL:   "https://example.com",
					ShortCode: "abcd0001",
				},
				{
					LongURL:   "https://github.com",
					ShortCode: "abcd0002",
				},
				{
					LongURL:   "https://google.com",
					ShortCode: "abcd0003",
				},
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := store.CreateURLs(ctx, tt.urls)
			assert.NoError(t, err)

			snapshot := store.CreateMemento()
			assert.Equal(t, tt.expected, len(snapshot.State))
		})
	}
}

func Test_InMemoryRepository_GetURLByShortCode(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryRepository()

	_, err := store.CreateURL(ctx, URL{
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
			longURL, found := store.GetURLByShortCode(ctx, tt.shortURL)

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

func Benchmark_InMemoryRepository_GetURLByShortCode(b *testing.B) {
	ctx := context.Background()
	store := NewInMemoryRepository()

	_, err := store.CreateURL(ctx, URL{
		LongURL:   "https://example.com",
		ShortCode: "abcd1234",
	})
	assert.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		store.GetURLByShortCode(ctx, "abcd1234")
	}
}

func Test_InMemoryRepository_GetURLsByUserID(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryRepository()

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
				_, err := store.CreateURL(ctx, URL{
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
				_, err := store.CreateURL(ctx, URL{
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
				_, err := store.CreateURL(ctx, URL{
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

			urls, total, err := store.GetURLsByUserID(ctx, tt.UserID, 25, 0)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.count, total)

			if tt.expected.count > 0 {
				assert.Equal(t, tt.expected.UUID, urls[0].UUID)
				assert.Equal(t, tt.expected.LongURL, urls[0].LongURL)
				assert.Equal(t, tt.expected.ShortCode, urls[0].ShortCode)
			}

			t.Cleanup(func() {
				store.Clear()
			})
		})
	}
}

func Test_InMemoryRepository_DeleteURLsByUserID(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryRepository()

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
		deleted       bool
		expected      int
	}{
		{
			name:          "Success",
			currentUserID: UserUUID1,
			ownerID:       UserUUID1,
			before: func() {
				_, err := store.CreateURL(ctx, URL{
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
			deleted:  true,
			expected: 2,
		},
		{
			name:          "Not owned",
			currentUserID: UserUUID1,
			ownerID:       UserUUID2,
			before: func() {
				_, err := store.CreateURL(ctx, URL{
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
			deleted:  false,
			expected: 1,
		},
		{
			name:          "Not Found",
			currentUserID: UserUUID1,
			ownerID:       UserUUID1,
			before: func() {
				_, err := store.CreateURL(ctx, URL{
					UUID:      UUID1,
					LongURL:   "https://google.com",
					ShortCode: "abcd0001",
					UserUUID:  UserUUID1,
				})
				assert.NoError(t, err)
			},
			params: db.DeleteURLsByUserIDAndShortCodesParams{
				UserUUID:   UserUUID1,
				ShortCodes: []string{"abcd0002"},
			},
			deleted:  false,
			expected: 1,
		},
		{
			name:          "Empty",
			currentUserID: UserUUID1,
			ownerID:       UserUUID1,
			before: func() {
				_, err := store.CreateURL(ctx, URL{
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
			deleted:  false,
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before()

			err := store.DeleteURLsByUserID(ctx, tt.params.UserUUID, tt.params.ShortCodes)
			assert.NoError(t, err)

			snapshot := store.CreateMemento()
			assert.Equal(t, tt.expected, len(snapshot.State))

			for _, url := range snapshot.State {
				if tt.deleted {
					assert.NotEmpty(t, url.DeletedAt)
				} else {
					assert.Empty(t, url.DeletedAt)
				}
			}

			t.Cleanup(func() {
				store.Clear()
			})
		})
	}
}

func Test_InMemoryRepository_CreateMemento(t *testing.T) {
	ctx := context.Background()

	UUID := uuid.MustParse("6455bd07-e431-4851-af3c-4f703f726639")

	type result struct {
		memento *Memento
		err     error
	}

	tests := []struct {
		name     string
		before   func(store InMemory)
		expected result
	}{
		{
			name: "Success",
			before: func(store InMemory) {
				store.CreateURL(ctx, URL{
					UUID:      UUID,
					LongURL:   "http://example.com",
					ShortCode: "abcd1234",
				})
			},
			expected: result{
				memento: &Memento{
					State: []URL{
						{
							UUID:      UUID,
							LongURL:   "http://example.com",
							ShortCode: "abcd1234",
						},
					},
				},
				err: nil,
			},
		},
		{
			name:   "Empty",
			before: func(_ InMemory) {},
			expected: result{
				memento: &Memento{
					State: []URL(nil),
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewInMemoryRepository()
			tt.before(store)

			memento := store.CreateMemento()
			assert.Equal(t, tt.expected.memento, memento)
		})
	}
}

func Test_InMemoryRepository_Restore(t *testing.T) {
	ctx := context.Background()

	UUID := uuid.MustParse("6455bd07-e431-4851-af3c-4f703f726639")

	type result struct {
		memento *Memento
		err     error
	}

	tests := []struct {
		name     string
		before   func(store InMemory)
		expected result
	}{
		{
			name:   "Success",
			before: func(_ InMemory) {},
			expected: result{
				memento: &Memento{
					State: []URL{
						{
							UUID:      UUID,
							LongURL:   "http://example.com",
							ShortCode: "abcd1234",
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "Empty",
			before: func(store InMemory) {
				store.CreateURL(ctx, URL{
					UUID:      UUID,
					LongURL:   "http://example.com",
					ShortCode: "abcd1234",
				})
			},
			expected: result{
				memento: &Memento{
					State: []URL(nil),
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := NewInMemoryRepository()
			tt.before(store)

			store.Restore(tt.expected.memento)
			memento := store.CreateMemento()
			assert.Equal(t, tt.expected.memento, memento)
		})
	}
}

func Test_InMemoryRepository_Clear(t *testing.T) {
	ctx := context.Background()
	store := NewInMemoryRepository()

	UUID := uuid.MustParse("6455bd07-e431-4851-af3c-4f703f726639")
	UserUUID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	tests := []struct {
		name     string
		before   func()
		expected int
	}{
		{
			name: "Success",
			before: func() {
				_, err := store.CreateURL(ctx, URL{
					UUID:      UUID,
					LongURL:   "https://example.com",
					ShortCode: "abcd1234",
					UserUUID:  UserUUID,
				})
				assert.NoError(t, err)
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store.Clear()

			snapshot := store.CreateMemento()
			assert.Equal(t, tt.expected, len(snapshot.State))
		})
	}
}
