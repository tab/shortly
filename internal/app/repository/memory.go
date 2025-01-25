package repository

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

// InMemory is an interface for in-memory storage
type InMemory interface {
	Repository
	CreateMemento() *Memento
	Restore(m *Memento)
	Clear()
}

// InMemoryRepo is a repository for in-memory storage
type InMemoryRepo struct {
	data sync.Map
}

// NewInMemoryRepository creates a new in-memory repository instance
func NewInMemoryRepository() InMemory {
	return &InMemoryRepo{}
}

// CreateURL creates a new URL record
func (m *InMemoryRepo) CreateURL(_ context.Context, url URL) (*URL, error) {
	m.data.Store(url.ShortCode, url)
	return &url, nil
}

// CreateURLs creates new URL records
func (m *InMemoryRepo) CreateURLs(_ context.Context, urls []URL) error {
	for _, url := range urls {
		m.data.Store(url.ShortCode, url)
	}
	return nil
}

// GetURLByShortCode returns a URL record by short code
func (m *InMemoryRepo) GetURLByShortCode(_ context.Context, shortCode string) (*URL, bool) {
	value, ok := m.data.Load(shortCode)
	if !ok {
		return nil, false
	}

	url, ok := value.(URL)
	if !ok {
		return nil, false
	}

	return &url, true
}

// GetURLsByUserID returns URL records by user ID
func (m *InMemoryRepo) GetURLsByUserID(_ context.Context, id uuid.UUID, limit, offset int64) ([]URL, int, error) {
	var results []URL

	m.data.Range(func(_, value interface{}) bool {
		url, ok := value.(URL)
		if ok && url.UserUUID == id && url.DeletedAt.IsZero() {
			results = append(results, url)
		}
		return true
	})

	total := len(results)

	start := int(offset)
	if start > total {
		start = total
	}

	end := int(offset + limit)
	if end > total {
		end = total
	}

	return results[start:end], total, nil
}

// DeleteURLsByUserID deletes URL records by user ID
func (m *InMemoryRepo) DeleteURLsByUserID(_ context.Context, id uuid.UUID, shortCodes []string) error {
	for _, shortCode := range shortCodes {
		value, ok := m.data.Load(shortCode)
		if !ok {
			continue
		}

		url, ok := value.(URL)
		if !ok {
			continue
		}

		if url.UserUUID == id && url.DeletedAt.IsZero() {
			url.DeletedAt = time.Now()
			m.data.Store(shortCode, url)
		}
	}

	return nil
}

// CreateMemento creates a memento of the current state
func (m *InMemoryRepo) CreateMemento() *Memento {
	var results []URL

	m.data.Range(func(_, value interface{}) bool {
		url, ok := value.(URL)
		if ok {
			results = append(results, url)
		}
		return true
	})

	return &Memento{State: results}
}

// Restore restores the state from a memento
func (m *InMemoryRepo) Restore(memento *Memento) {
	m.data = sync.Map{}

	for _, url := range memento.State {
		m.data.Store(url.ShortCode, url)
	}
}

// Clear clears the repository
func (m *InMemoryRepo) Clear() {
	m.data = sync.Map{}
}
