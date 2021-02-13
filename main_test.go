package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPut(t *testing.T) {
	store := make(store)
	application := NewApplication(store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/key", strings.NewReader("value"))
	application.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "value", w.Body.String())
	assert.Contains(t, store, "key")
	assert.Equal(t, []byte("value"), store["key"])
}

func TestGet(t *testing.T) {
	store := make(store)
	store["key"] = []byte("value")
	application := NewApplication(store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/key", nil)
	application.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "value", w.Body.String())
}

func TestGetMissingKey(t *testing.T) {
	store := make(store)
	application := NewApplication(store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/key", nil)
	application.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}

func TestDelete(t *testing.T) {
	store := make(store)
	store["key"] = []byte("value")
	application := NewApplication(store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/key", nil)
	application.ServeHTTP(w, req)

	assert.Equal(t, 204, w.Code)
	assert.NotContains(t, store, "key")
}

func TestGetIndex(t *testing.T) {
	store := make(store)
	store["key2"] = []byte("value2")
	store["key1"] = []byte("value1")
	application := NewApplication(store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	application.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "key1\nkey2", w.Body.String())
}
