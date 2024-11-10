package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
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
			err := store.CreateURL(ctx, tt.url)
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

	err := store.CreateURL(ctx, URL{
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

func Test_InMemoryRepository_CreateMemento(t *testing.T) {
	ctx := context.Background()

	UUID, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f726639")

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

	UUID, _ := uuid.Parse("6455bd07-e431-4851-af3c-4f703f726639")

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
