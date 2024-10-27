package repository

import "sync"

type InMemoryRepository struct {
	data map[string]URL
	mu   sync.RWMutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data: make(map[string]URL),
	}
}

func (store *InMemoryRepository) Set(url URL) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.data[url.ShortCode] = url

	return nil
}

func (store *InMemoryRepository) Get(shortCode string) (*URL, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	url, found := store.data[shortCode]
	if !found {
		return nil, false
	}
	return &url, true
}

func (store *InMemoryRepository) GetAll() []URL {
	store.mu.RLock()
	defer store.mu.RUnlock()

	results := make([]URL, 0, len(store.data))
	for _, v := range store.data {
		results = append(results, v)
	}
	return results
}
