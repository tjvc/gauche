package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tjvc/gauche/internal/logging"
	"github.com/tjvc/gauche/internal/store"
)

func TestPut(t *testing.T) {
	store := store.New()
	server := buildServer(&store)
	defer server.Close()
	url := fmt.Sprintf("%s/key", server.URL)
	req, _ := http.NewRequest("PUT", url, strings.NewReader("value"))

	response, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 200, response.StatusCode)
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, []byte("value"), body)
	value, _ := store.Get("key")
	assert.Equal(t, []byte("value"), value)
}

func TestGet(t *testing.T) {
	store := store.New()
	store.Set("key", []byte("value"))
	server := buildServer(&store)
	defer server.Close()
	url := fmt.Sprintf("%s/key", server.URL)
	req, _ := http.NewRequest("GET", url, nil)

	response, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 200, response.StatusCode)
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, []byte("value"), body)
}

func TestGetMissingKey(t *testing.T) {
	store := store.New()
	server := buildServer(&store)
	defer server.Close()
	url := fmt.Sprintf("%s/key", server.URL)
	req, _ := http.NewRequest("GET", url, nil)

	response, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 404, response.StatusCode)
}

func TestDelete(t *testing.T) {
	store := store.New()
	store.Set("key", []byte("value"))
	server := buildServer(&store)
	defer server.Close()
	url := fmt.Sprintf("%s/key", server.URL)
	req, _ := http.NewRequest("DELETE", url, nil)

	response, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 204, response.StatusCode)
	_, present := store.Get("key")
	assert.False(t, present)
}

func TestGetIndex(t *testing.T) {
	store := store.New()
	store.Set("key2", []byte("value2"))
	store.Set("key1", []byte("value1"))
	server := buildServer(&store)
	defer server.Close()
	req, _ := http.NewRequest("GET", server.URL, nil)

	response, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 200, response.StatusCode)
	body, _ := ioutil.ReadAll(response.Body)
	assert.Equal(t, []byte("key1\nkey2"), body)
}

func TestInvalidMethod(t *testing.T) {
	store := store.New()
	server := buildServer(&store)
	defer server.Close()
	url := fmt.Sprintf("%s/key", server.URL)
	req, _ := http.NewRequest("POST", url, nil)

	response, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 405, response.StatusCode)
}

func TestPutMissingKey(t *testing.T) {
	store := store.New()
	server := buildServer(&store)
	defer server.Close()
	req, _ := http.NewRequest("PUT", server.URL, nil)

	response, _ := http.DefaultClient.Do(req)

	assert.Equal(t, 405, response.StatusCode)
}

func buildServer(store *store.Store) *httptest.Server {
	application := buildApplication(store)
	server := httptest.NewServer(mainHandler(application))
	return server
}

func buildApplication(store *store.Store) application {
	logger := nullLogger{}

	return application{
		store:  store,
		logger: logger,
	}
}

type nullLogger struct{}

func (nullLogger) Write(logging.LogEntry) {
}
