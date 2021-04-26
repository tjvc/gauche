package main

import (
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
)

type store struct {
	store map[string][]byte
}

// Application wraps a data store and HTTP handler
type Application struct {
	store       *store
	httpHandler *gin.Engine
}

func (application *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	application.httpHandler.ServeHTTP(w, r)
}

// Run starts the application
func (application *Application) Run() {
	application.httpHandler.Run()
}

func putHandler(store *store) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Params.ByName("key")
		value, _ := ioutil.ReadAll(c.Request.Body)
		store.store[key] = value
		c.Data(200, "text/plain", value)
	}
}

func getHandler(store *store) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Params.ByName("key")
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
		delete(store.store, key)
		c.Status(204)
	}
}

func getIndexHandler(store *store) func(*gin.Context) {
	return func(c *gin.Context) {
		keys := make([]string, len(store.store))
		i := 0
		for key := range store.store {
			keys[i] = key
			i++
		}
		sort.Strings(keys)
		c.String(200, strings.Join(keys, "\n"))
	}
}

// NewApplication returns a new Application
func NewApplication(store *store) Application {
	router := gin.Default()

	router.PUT("/:key", putHandler(store))
	router.GET("/:key", getHandler(store))
	router.DELETE("/:key", deleteHandler(store))
	router.GET("/", getIndexHandler(store))

	application := Application{
		store:       store,
		httpHandler: router,
	}

	return application
}

func main() {
	store := store{
		store: make(map[string][]byte),
	}
	application := NewApplication(&store)
	application.Run()
}
