package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreSet(t *testing.T) {
	store := newStore()
	store.set("key", []byte("value"))

	value, _ := store.get("key")
	assert.Equal(t, []byte("value"), value)
}

func TestStoreGet(t *testing.T) {
	store := newStore()
	store.set("key", []byte("value"))

	value, present := store.get("key")
	assert.Equal(t, []byte("value"), value)
	assert.True(t, present)
}

func TestStoreGetMissingKey(t *testing.T) {
	store := newStore()

	value, present := store.get("key")
	assert.Nil(t, value)
	assert.False(t, present)
}

func TestStoreDelete(t *testing.T) {
	store := newStore()
	store.set("key", []byte("value"))

	store.delete("key")
	value, present := store.get("key")
	assert.Nil(t, value)
	assert.False(t, present)
}

func TestStoreGetKeys(t *testing.T) {
	store := newStore()
	store.set("key2", []byte("value2"))
	store.set("key1", []byte("value1"))

	keys := store.keys()
	assert.Equal(t, []string{"key1", "key2"}, keys)
}
