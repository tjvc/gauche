package store

import (
	"sort"
	"sync"
)

type Store struct {
	sync.RWMutex
	store map[string][]byte
	Size  int
}

func New() Store {
	return Store{
		store: make(map[string][]byte),
	}
}

func (store *Store) Get(key string) ([]byte, bool) {
	store.RLock()
	defer store.RUnlock()
	value, present := store.store[key]
	return value, present
}

func (store *Store) Set(key string, value []byte) {
	store.Lock()
	defer store.Unlock()
	store.store[key] = value
	store.Size += len(key)
	store.Size += len(value)
}

func (store *Store) Delete(key string) {
	store.Lock()
	defer store.Unlock()
	delete(store.store, key)
}

func (store *Store) Keys() []string {
	keys := make([]string, len(store.store))
	i := 0
	store.RLock()
	defer store.RUnlock()
	for key := range store.store {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}
