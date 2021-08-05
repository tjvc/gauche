package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPut(t *testing.T) {
	store := newStore()
	application := newApplication(&store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/key", strings.NewReader("value"))
	application.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "value", w.Body.String())
	assert.Contains(t, store.store, "key")
	assert.Equal(t, []byte("value"), store.store["key"])
}

func TestGet(t *testing.T) {
	store := newStore()
	store.set("key", []byte("value"))
	application := newApplication(&store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/key", nil)
	application.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "value", w.Body.String())
}

func TestGetMissingKey(t *testing.T) {
	store := newStore()
	application := newApplication(&store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/key", nil)
	application.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestDelete(t *testing.T) {
	store := newStore()
	store.set("key", []byte("value"))
	application := newApplication(&store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/key", nil)
	application.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
	assert.NotContains(t, store.store, "key")
}

func TestGetIndex(t *testing.T) {
	store := newStore()
	store.set("key2", []byte("value2"))
	store.set("key1", []byte("value1"))
	application := newApplication(&store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	application.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "key1\nkey2", w.Body.String())
}
