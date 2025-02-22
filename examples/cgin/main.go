package main

import (
	"cgin"
	"net/http"
)

func main() {
	r := cgin.Default()
	r.GET("/", func(c *cgin.Context) {
		c.String(http.StatusOK, "Hello")
	})
	// index out of range for testing Recovery()
	r.GET("/panic", func(c *cgin.Context) {
		names := []string{""}
		c.String(http.StatusOK, names[100])
	})

	r.Run(":8888")
}
