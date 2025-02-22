package cgin

import (
	"strings"
	"sync"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    map[string]*node{},
		handlers: map[string]HandlerFunc{},
	}
}

func parsePattern(pattern string) []string {
	strs := strings.Split(pattern, "/")

	parts := []string{}
	for _, item := range strs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}

	return parts
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}

	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

var mapPool = sync.Pool{
	New: func() interface{} {
		return make(map[string]string)
	},
}

func (r *router) getRoute(method, pattern string) (*node, map[string]string) {
	searchParts := parsePattern(pattern)
	params := mapPool.Get().(map[string]string)
	// params := map[string]string{}
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}

		return n, params
	}

	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]

	if ok {
		return root.getChildren()
	}

	return nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(404, "404 NOT FOUND: %s\n", c.Path)
		})
	}

	c.Next()
}
