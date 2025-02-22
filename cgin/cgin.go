package cgin

import (
	"log"
	"net/http"
	"path"
	"strings"
	"sync"
)

type HandlerFunc func(*Context)

type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc
		parent      *RouterGroup
		engine      *Engine
	}

	Engine struct {
		RouterGroup
		router *router
		groups []*RouterGroup
	}
)

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{&engine.RouterGroup}
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method, path string, handler HandlerFunc) {
	pattern := group.prefix + path
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(path string, handler HandlerFunc) {
	group.addRoute("GET", path, handler)
}

func (group *RouterGroup) POST(path string, handler HandlerFunc) {
	group.addRoute("POST", path, handler)
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}

		fileServer.ServeHTTP(c.ResponseWriter, c.Request)
	}
}

func (group *RouterGroup) Static(relativePath string, root string) {
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
}

// Regarding HTML templates, they are not implemented here.
// This is not the primary focus of the project.
// HTML templates are used infrequently in this context.

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

var MiddlewaresPool = sync.Pool{
	New: func() interface{} {
		return make([]HandlerFunc, 0, 2)
	},
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	// c := newContext(w, r)
	c := ContextPool.Get().(*Context)
	initContext(c, w, r)
	c.handlers = middlewares
	c.engine = engine
	engine.router.handle(c)
	for k := range c.Params {
		delete(c.Params, k)
	}
	mapPool.Put(c.Params)
	ContextPool.Put(c)
}
