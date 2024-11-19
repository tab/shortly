package repository

import (
	"context"
	"sync"

	"github.com/google/uuid"
)

type InMemory interface {
	Repository
	CreateMemento() *Memento
	Restore(m *Memento)
	Clear()
}

type InMemoryRepo struct {
	data sync.Map
}

func NewInMemoryRepository() InMemory {
	return &InMemoryRepo{}
}

func (m *InMemoryRepo) CreateURL(_ context.Context, url URL) (*URL, error) {
	m.data.Store(url.ShortCode, url)
	return &url, nil
}

func (m *InMemoryRepo) CreateURLs(_ context.Context, urls []URL) error {
	for _, url := range urls {
		m.data.Store(url.ShortCode, url)
	}
	return nil
}

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

func (m *InMemoryRepo) GetURLsByUserID(_ context.Context, id uuid.UUID) ([]URL, error) {
	var results []URL

	m.data.Range(func(_, value interface{}) bool {
		url, ok := value.(URL)
		if ok && url.UserUUID == id {
			results = append(results, url)
		}
		return true
	})

	return results, nil
}

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

func (m *InMemoryRepo) Restore(memento *Memento) {
	m.data = sync.Map{}

	for _, url := range memento.State {
		m.data.Store(url.ShortCode, url)
	}
}

func (m *InMemoryRepo) Clear() {
	m.data = sync.Map{}
}
