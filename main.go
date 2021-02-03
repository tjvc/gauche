package main

import "io/ioutil"
import "net/http"
import "github.com/gin-gonic/gin"

type store map[string]string

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
		b, _ := ioutil.ReadAll(c.Request.Body)
		value := string(b)
		store[key] = value
		c.String(200, value)
	})

	router.GET("/:key", func(c *gin.Context) {
		key := c.Params.ByName("key")
		if value, present := store[key]; present {
			c.String(200, value)
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
