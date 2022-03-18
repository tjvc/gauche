package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type application struct {
	store  *store
	logger logger
	mux    *http.ServeMux
}

func (application *application) run() {
	http.ListenAndServe(":8080", application.mux)
}

func putHandler(w http.ResponseWriter, key string, r *http.Request, store *store) {
	value, _ := ioutil.ReadAll(r.Body)
	store.set(key, value)
	w.WriteHeader(http.StatusOK)
	w.Write(value)
}

func getHandler(w http.ResponseWriter, key string, store *store) {
	if value, present := store.get(key); present {
		w.WriteHeader(http.StatusOK)
		w.Write(value)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func deleteHandler(w http.ResponseWriter, key string, store *store) {
	store.delete(key)
	w.WriteHeader(http.StatusNoContent)
}

func getIndexHandler(w http.ResponseWriter, store *store) {
	keys := store.keys()
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, strings.Join(keys, "\n"))
}

func recoveryMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		handler.ServeHTTP(w, r)
	})
}

func (application application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		if r.Method == http.MethodGet {
			getIndexHandler(w, application.store)
			return
		}

		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Path[1:]

	switch r.Method {
	case http.MethodGet:
		getHandler(w, key, application.store)
	case http.MethodPut:
		putHandler(w, key, r, application.store)
	case http.MethodDelete:
		deleteHandler(w, key, application.store)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func newApplication(store *store, logger logger) application {
	mux := http.NewServeMux()

	application := application{
		store:  store,
		logger: logger,
		mux:    mux,
	}

	mux.Handle("/", loggingMiddleware(logger, recoveryMiddleware(application)))

	return application
}

func main() {
	store := newStore()
	logger := jsonLogger{}
	application := newApplication(&store, logger)
	application.run()
}
