package main

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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
		store.set(key, value)
		c.Data(200, "text/plain", value)
	}
}

func getHandler(store *store) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Params.ByName("key")
		if value, present := store.get(key); present {
			c.Data(200, "text/plain", value)
		} else {
			c.Status(404)
		}
	}
}

func deleteHandler(store *store) func(*gin.Context) {
	return func(c *gin.Context) {
		key := c.Params.ByName("key")
		store.delete(key)
		c.Status(204)
	}
}

func getIndexHandler(store *store) func(*gin.Context) {
	return func(c *gin.Context) {
		keys := store.keys()
		c.String(200, strings.Join(keys, "\n"))
	}
}

func recoveryMiddleware(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			c.AbortWithStatus(500)
		}
	}()
	c.Next()
}

func newApplication(store *store, logger logger) application {
	gin.SetMode("release")
	router := gin.New()

	router.Use(loggingMiddleware(logger))
	router.Use(recoveryMiddleware)

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
	store := newStore()
	logger := jsonLogger{}
	application := newApplication(&store, logger)
	application.run()
}
