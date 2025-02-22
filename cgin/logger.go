package cgin

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()

		c.Next()

		log.Printf("[%d] %s costs %v", 200 /*c.StatusCode*/, c.Request.RequestURI, time.Since(t))
	}
}
