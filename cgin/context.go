package cgin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type H map[string]interface{}

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request

	Path   string
	Method string
	Params map[string]string

	StatusCode int
	handlers   []HandlerFunc
	index      int

	engine *Engine
}

var ContextPool = sync.Pool{
	New: func() interface{} {
		return &Context{}
	},
}

func initContext(c *Context, w http.ResponseWriter, r *http.Request) {
	c.ResponseWriter = w
	c.Request = r
	c.Path = r.URL.Path
	c.Method = r.Method
	c.index = -1
}

func (c *Context) reset() {
	c.ResponseWriter = nil
	c.Request = nil
	c.Path = ""
	c.Method = ""
	c.Params = nil

	c.handlers = nil
	c.StatusCode = 0
	c.index = -1

	c.engine = nil
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		ResponseWriter: w,
		Request:        r,
		Path:           r.URL.Path,
		Method:         r.Method,
		index:          -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Fail(code int, err string) {
	c.index = len(c.handlers)
	c.JSON(code, H{"message": err})
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.ResponseWriter.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.ResponseWriter.Header().Set(key, value)
}

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.ResponseWriter.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.ResponseWriter)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.ResponseWriter, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.ResponseWriter.Write(data)
}

// HTML template render, we will implement it later
// func (c *Context) HTML(code int, name string, data interface{}) {}
