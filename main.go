package main

import "io/ioutil"
import "github.com/gin-gonic/gin"

type store map[string]string

func setupRouter(store store) *gin.Engine {
	r := gin.Default()

	r.PUT("/:key", func(c *gin.Context) {
		key := c.Params.ByName("key")
		b, _ := ioutil.ReadAll(c.Request.Body)
		value := string(b)
		store[key] = value
		c.String(200, value)
	})

	r.GET("/:key", func(c *gin.Context) {
		key := c.Params.ByName("key")
		if value, present := store[key]; present {
			c.String(200, value)
		} else {
			c.Status(404)
		}
	})

	return r
}

func main() {
	store := make(store)
	r := setupRouter(store)
	r.Run(":8080")
}
