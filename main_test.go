package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPut(t *testing.T) {
	store := newStore()
	server := buildServer(&store)
	defer server.Close()
	url := fmt.Sprintf("%s/key", server.URL)
	req, _ := http.NewRequest("PUT", url, strings.NewReader("value"))

	response, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 200, response.StatusCode)
	buf := new(strings.Builder)
	io.Copy(buf, response.Body)
	assert.Equal(t, "value", buf.String())
	assert.Contains(t, store.store, "key")
	assert.Equal(t, []byte("value"), store.store["key"])
}

func TestGet(t *testing.T) {
	store := newStore()
	store.set("key", []byte("value"))
	server := buildServer(&store)
	defer server.Close()
	url := fmt.Sprintf("%s/key", server.URL)
	req, _ := http.NewRequest("GET", url, nil)

	response, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 200, response.StatusCode)
	buf := new(strings.Builder)
	io.Copy(buf, response.Body)
	assert.Equal(t, "value", buf.String())
}

func TestGetMissingKey(t *testing.T) {
	store := newStore()
	server := buildServer(&store)
	defer server.Close()
	url := fmt.Sprintf("%s/key", server.URL)
	req, _ := http.NewRequest("GET", url, nil)

	response, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 404, response.StatusCode)
}

func TestDelete(t *testing.T) {
	store := newStore()
	store.set("key", []byte("value"))
	server := buildServer(&store)
	defer server.Close()
	url := fmt.Sprintf("%s/key", server.URL)
	req, _ := http.NewRequest("DELETE", url, nil)

	response, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 204, response.StatusCode)
	assert.NotContains(t, store.store, "key")
}

func TestGetIndex(t *testing.T) {
	store := newStore()
	store.set("key2", []byte("value2"))
	store.set("key1", []byte("value1"))
	server := buildServer(&store)
	defer server.Close()
	req, _ := http.NewRequest("GET", server.URL, nil)

	response, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 200, response.StatusCode)
	buf := new(strings.Builder)
	io.Copy(buf, response.Body)
	assert.Equal(t, "key1\nkey2", buf.String())
}

func TestInvalidMethod(t *testing.T) {
	store := newStore()
	server := buildServer(&store)
	defer server.Close()
	url := fmt.Sprintf("%s/key", server.URL)
	req, _ := http.NewRequest("POST", url, nil)

	response, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 405, response.StatusCode)
}

func TestPutMissingKey(t *testing.T) {
	store := newStore()
	server := buildServer(&store)
	defer server.Close()
	req, _ := http.NewRequest("PUT", server.URL, nil)

	response, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 405, response.StatusCode)
}

func buildServer(store *store) *httptest.Server {
	application := buildApplication(store)
	server := httptest.NewServer(application.mux)
	return server
}

func buildApplication(store *store) application {
	logger := nullLogger{}
	return newApplication(store, logger)
}

type nullLogger struct{}

func (nullLogger) write(logEntry) {
}
