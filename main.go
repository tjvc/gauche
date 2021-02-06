package main

import "io/ioutil"
import "net/http"
import "github.com/gin-gonic/gin"

type store map[string][]byte

type application struct {
	store       store
	httpHandler *gin.Engine
}

func (application *application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	application.httpHandler.ServeHTTP(w, r)
}

func (application *application) Run() {
	application.httpHandler.Run()
}

func newApplication() application {
	store := make(store)
	router := gin.Default()

	router.PUT("/:key", func(c *gin.Context) {
		key := c.Params.ByName("key")
		value, _ := ioutil.ReadAll(c.Request.Body)
		store[key] = value
		c.Data(200, "text/plain", value)
	})

	router.GET("/:key", func(c *gin.Context) {
		key := c.Params.ByName("key")
		if value, present := store[key]; present {
			c.Data(200, "text/plain", value)
		} else {
			c.Status(404)
		}
	})

	application := application{
		store:       store,
		httpHandler: router,
	}

	return application
}

func main() {
	application := newApplication()
	application.Run()
}
