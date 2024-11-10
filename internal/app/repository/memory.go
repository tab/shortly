package repository

import (
	"context"
	"sync"
)

type InMemory interface {
	Repository
	CreateMemento() *Memento
	Restore(m *Memento)
}

type InMemoryRepo struct {
	data sync.Map
}

func NewInMemoryRepository() InMemory {
	return &InMemoryRepo{}
}

func (m *InMemoryRepo) CreateURL(_ context.Context, url URL) error {
	m.data.Store(url.ShortCode, url)
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
