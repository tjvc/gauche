package main

import "sync"

type store struct {
	sync.RWMutex
	store map[string][]byte
}

func newStore() store {
	return store{
		store: make(map[string][]byte),
	}
}
