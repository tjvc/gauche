package main

import "io/ioutil"
import "net/http"
import "github.com/gin-gonic/gin"

type store map[string][]byte

// Application wraps a data store and HTTP handler
type Application struct {
	store       store
	httpHandler *gin.Engine
}

func (application *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	application.httpHandler.ServeHTTP(w, r)
}

// Run starts the application
func (application *Application) Run() {
	application.httpHandler.Run()
}

func putHandler(store store) func(*gin.Context) {
	fn := func(c *gin.Context) {
		key := c.Params.ByName("key")
		value, _ := ioutil.ReadAll(c.Request.Body)
		store[key] = value
		c.Data(200, "text/plain", value)
	}
	return fn
}

func getHandler(store store) func(*gin.Context) {
	fn := func(c *gin.Context) {
		key := c.Params.ByName("key")
		if value, present := store[key]; present {
			c.Data(200, "text/plain", value)
		} else {
			c.Status(404)
		}
	}
	return fn
}

// NewApplication returns a new Application
func NewApplication() Application {
	store := make(store)
	router := gin.Default()

	router.PUT("/:key", putHandler(store))
	router.GET("/:key", getHandler(store))

	application := Application{
		store:       store,
		httpHandler: router,
	}

	return application
}

func main() {
	application := NewApplication()
	application.Run()
}
