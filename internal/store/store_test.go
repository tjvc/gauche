package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreSet(t *testing.T) {
	store := New()
	store.Set("key", []byte("value"))

	value, _ := store.Get("key")
	assert.Equal(t, []byte("value"), value)
}

func TestStoreGet(t *testing.T) {
	store := New()
	store.Set("key", []byte("value"))

	value, present := store.Get("key")
	assert.Equal(t, []byte("value"), value)
	assert.True(t, present)
}

func TestStoreGetMissingKey(t *testing.T) {
	store := New()

	value, present := store.Get("key")
	assert.Nil(t, value)
	assert.False(t, present)
}

func TestStoreDelete(t *testing.T) {
	store := New()
	store.Set("key", []byte("value"))

	store.Delete("key")
	value, present := store.Get("key")
	assert.Nil(t, value)
	assert.False(t, present)
}

func TestStoreGetKeys(t *testing.T) {
	store := New()
	store.Set("key2", []byte("value2"))
	store.Set("key1", []byte("value1"))

	keys := store.Keys()
	assert.Equal(t, []string{"key1", "key2"}, keys)
}
