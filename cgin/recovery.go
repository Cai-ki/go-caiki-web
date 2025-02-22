package cgin

import (
	"log"
	"net/http"
	"runtime/debug"
)

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				debug.PrintStack()
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}
