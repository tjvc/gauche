package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

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

	if response.StatusCode != 200 {
		t.Errorf("got %d, want %d", response.StatusCode, 200)
	}
	want := []byte("value")
	body, _ := ioutil.ReadAll(response.Body)
	res := bytes.Compare(body, want)
	if res != 0 {
		t.Errorf("got %s, want %s", body, want)
	}
	value, _ := store.Get("key")
	res = bytes.Compare(value, want)
	if res != 0 {
		t.Errorf("got %s, want %s", value, want)
	}
}

func TestGet(t *testing.T) {
	store := store.New()
	store.Set("key", []byte("value"))
	server := buildServer(&store)
	defer server.Close()
	url := fmt.Sprintf("%s/key", server.URL)
	req, _ := http.NewRequest("GET", url, nil)

	response, _ := http.DefaultClient.Do(req)

	if response.StatusCode != 200 {
		t.Errorf("got %d, want %d", response.StatusCode, 200)
	}
	want := []byte("value")
	body, _ := ioutil.ReadAll(response.Body)
	res := bytes.Compare(body, want)
	if res != 0 {
		t.Errorf("got %s, want %s", body, want)
	}
}

func TestGetMissingKey(t *testing.T) {
	store := store.New()
	server := buildServer(&store)
	defer server.Close()
	url := fmt.Sprintf("%s/key", server.URL)
	req, _ := http.NewRequest("GET", url, nil)

	response, _ := http.DefaultClient.Do(req)

	if response.StatusCode != 404 {
		t.Errorf("got %d, want %d", response.StatusCode, 404)
	}
}

func TestDelete(t *testing.T) {
	store := store.New()
	store.Set("key", []byte("value"))
	server := buildServer(&store)
	defer server.Close()
	url := fmt.Sprintf("%s/key", server.URL)
	req, _ := http.NewRequest("DELETE", url, nil)

	response, _ := http.DefaultClient.Do(req)

	if response.StatusCode != 204 {
		t.Errorf("got %d, want %d", response.StatusCode, 204)
	}
	_, present := store.Get("key")
	if present != false {
		t.Error("got true, want false")
	}
}

func TestGetIndex(t *testing.T) {
	store := store.New()
	store.Set("key2", []byte("value2"))
	store.Set("key1", []byte("value1"))
	server := buildServer(&store)
	defer server.Close()
	req, _ := http.NewRequest("GET", server.URL, nil)

	response, _ := http.DefaultClient.Do(req)

	if response.StatusCode != 200 {
		t.Errorf("got %d, want %d", response.StatusCode, 200)
	}
	want := []byte("key1\nkey2")
	body, _ := ioutil.ReadAll(response.Body)
	res := bytes.Compare(body, want)
	if res != 0 {
		t.Errorf("got %s, want %s", body, want)
	}
}

func TestInvalidMethod(t *testing.T) {
	store := store.New()
	server := buildServer(&store)
	defer server.Close()
	url := fmt.Sprintf("%s/key", server.URL)
	req, _ := http.NewRequest("POST", url, nil)

	response, _ := http.DefaultClient.Do(req)

	if response.StatusCode != 405 {
		t.Errorf("got %d, want %d", response.StatusCode, 405)
	}
}

func TestPutMissingKey(t *testing.T) {
	store := store.New()
	server := buildServer(&store)
	defer server.Close()
	req, _ := http.NewRequest("PUT", server.URL, nil)

	response, _ := http.DefaultClient.Do(req)

	if response.StatusCode != 405 {
		t.Errorf("got %d, want %d", response.StatusCode, 405)
	}
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
