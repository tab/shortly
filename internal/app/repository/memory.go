package repository

import "sync"

type InMemoryRepository struct {
	data map[string]URL
	sync.RWMutex
}

func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		data: make(map[string]URL),
	}
}

func (store *InMemoryRepository) Set(url URL) {
	store.Lock()
	defer store.Unlock()
	store.data[url.ShortCode] = url
}

func (store *InMemoryRepository) Get(shortCode string) (*URL, bool) {
	store.RLock()
	defer store.RUnlock()

	url, found := store.data[shortCode]
	if !found {
		return nil, false
	}
	return &url, true
}
