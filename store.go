package main

import "sync"

type store struct {
	sync.RWMutex
	store map[string][]byte
}