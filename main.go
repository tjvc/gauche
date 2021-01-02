package main

import "io/ioutil"
import "github.com/gin-gonic/gin"

func setupRouter(store map[string]string) *gin.Engine {
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
		value := store[key]
		c.String(200, value)
	})

	return r
}

func main() {
	store := make(map[string]string)
	r := setupRouter(store)
	r.Run(":8080")
}
