package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPutGet(t *testing.T) {
	application := newApplication()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/key", strings.NewReader("value"))
	application.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "value", w.Body.String())

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/key", nil)
	application.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "value", w.Body.String())
}

func TestGet404(t *testing.T) {
	application := newApplication()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/key", nil)
	application.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
}
