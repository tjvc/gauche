package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tjvc/gauche/internal/store"
)

func Put(w http.ResponseWriter, key string, r *http.Request, store *store.Store) {
	value, err := io.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}

	if len(value) < 1 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	store.Set(key, value)
	w.WriteHeader(http.StatusOK)
	w.Write(value)
}

func Get(w http.ResponseWriter, key string, store *store.Store) {
	if value, present := store.Get(key); present {
		w.WriteHeader(http.StatusOK)
		w.Write(value)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func Delete(w http.ResponseWriter, key string, store *store.Store) {
	if _, present := store.Get(key); present {
		store.Delete(key)
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func Index(w http.ResponseWriter, store *store.Store) {
	keys := store.Keys()
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, strings.Join(keys, "\n"))
}
