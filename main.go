package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type application struct {
	store *store
}

func (application *application) run() {
	http.ListenAndServe(":8080", application)
}

func putHandler(w http.ResponseWriter, r *http.Request, store *store) {
	key := "key"
	value, _ := ioutil.ReadAll(r.Body)
	store.set(key, value)
	w.Write(value)
}

func getHandler(w http.ResponseWriter, r *http.Request, store *store) {
	key := "key"
	if value, present := store.get(key); present {
		w.Write(value)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func deleteHandler(w http.ResponseWriter, r *http.Request, store *store) {
	key := "key"
	store.delete(key)
	w.WriteHeader(http.StatusNoContent)
}

func getIndexHandler(w http.ResponseWriter, r *http.Request, store *store) {
	keys := store.keys()
	fmt.Fprint(w, strings.Join(keys, "\n"))
}

// func recoveryMiddleware(c *gin.Context) {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			c.AbortWithStatus(500)
// 		}
// 	}()
// 	c.Next()
// }

func (application application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		matched, _ := regexp.MatchString(`/[^/]+`, r.URL.Path)
		if r.URL.Path == "/" {
			getIndexHandler(w, r, application.store)
		} else if matched {
			getHandler(w, r, application.store)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodPut:
		putHandler(w, r, application.store)
	case http.MethodDelete:
		deleteHandler(w, r, application.store)
	}
}

func newApplication(store *store, logger logger) application {
	// router.Use(loggingMiddleware(logger))
	// router.Use(recoveryMiddleware)

	application := application{
		store: store,
	}

	mux := http.NewServeMux()
	mux.Handle("/key", application)

	return application
}

func main() {
	store := newStore()
	logger := jsonLogger{}
	application := newApplication(&store, logger)
	application.run()
}
