package main

import (
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type store struct {
	sync.RWMutex
	store map[string][]byte
}

type application struct {
	store       *store
	httpHandler *gin.Engine
}

func (application *application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	application.httpHandler.ServeHTTP(w, r)
}

func (application *application) run() {
	application.httpHandler.Run()
}

func putHandler(store *store) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Params.ByName("key")
		value, _ := ioutil.ReadAll(c.Request.Body)
		store.Lock()
		defer store.Unlock()
		store.store[key] = value
		c.Data(200, "text/plain", value)
	}
}

func getHandler(store *store) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Params.ByName("key")
		store.RLock()
		defer store.RUnlock()
		if value, present := store.store[key]; present {
			c.Data(200, "text/plain", value)
		} else {
			c.Status(404)
		}
	}
}

func deleteHandler(store *store) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Params.ByName("key")
		store.Lock()
		defer store.Unlock()
		delete(store.store, key)
		c.Status(204)
	}
}

func getIndexHandler(store *store) func(*gin.Context) {
	return func(c *gin.Context) {
		keys := make([]string, len(store.store))
		i := 0
		store.RLock()
		defer store.RUnlock()
		for key := range store.store {
			keys[i] = key
			i++
		}
		sort.Strings(keys)
		c.String(200, strings.Join(keys, "\n"))
	}
}

func newApplication(store *store) application {
	router := gin.Default()

	router.PUT("/:key", putHandler(store))
	router.GET("/:key", getHandler(store))
	router.DELETE("/:key", deleteHandler(store))
	router.GET("/", getIndexHandler(store))

	application := application{
		store:       store,
		httpHandler: router,
	}

	return application
}

func main() {
	store := store{
		store: make(map[string][]byte),
	}
	application := newApplication(&store)
	application.run()
}
