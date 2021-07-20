package main

import (
	"sort"
	"sync"
)

type store struct {
	sync.RWMutex
	store map[string][]byte
}

func newStore() store {
	return store{
		store: make(map[string][]byte),
	}
}

func (store *store) get(key string) ([]byte, bool) {
	store.RLock()
	defer store.RUnlock()
	value, present := store.store[key]
	return value, present
}

func (store *store) set(key string, value []byte) {
	store.Lock()
	defer store.Unlock()
	store.store[key] = value
}

func (store *store) delete(key string) {
	store.Lock()
	defer store.Unlock()
	delete(store.store, key)
}

func (store *store) keys() []string {
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
