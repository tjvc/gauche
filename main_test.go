package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetGet(t *testing.T) {
	store := make(store)
	router := setupRouter(store)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/key", strings.NewReader("value"))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "value", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/key", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "value", w.Body.String())
}
