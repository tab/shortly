package store

import "sync"

type URLStore struct {
	m map[string]string
	sync.RWMutex
}

func NewURLStore() *URLStore {
	return &URLStore{
		m: make(map[string]string),
	}
}

func (store *URLStore) Set(shortCode, longURL string) {
	store.Lock()
	defer store.Unlock()
	store.m[shortCode] = longURL
}

func (store *URLStore) Get(shortCode string) (string, bool) {
	store.RLock()
	defer store.RUnlock()
	longURL, found := store.m[shortCode]

	return longURL, found
}
