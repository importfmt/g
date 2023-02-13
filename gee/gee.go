package gee

import (
	"log"
	"net/http"
	"strings"
)


type HandlerFunc func(*Context)


type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}


type RouterGroup struct {
	prefix string
	middlewares []HandlerFunc
	parent *RouterGroup
	engine *Engine
}


func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}


func (group *RouterGroup) Group(prefix string) *RouterGroup{
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}


func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}


func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}


func (group *RouterGroup) PATCH(pattern string, handler HandlerFunc) {
	group.addRoute("PATCH", pattern, handler)
}


func (group *RouterGroup) DELETE(pattern string, handler HandlerFunc) {
	group.addRoute("DELETE", pattern, handler)
}

func (e *Engine) Run(addr string) (err error){
	return http.ListenAndServe(addr, e)
}


// func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
// 	e.router.addRoute(method, pattern, handler)
// }


// func (e *Engine) GET(pattern string, handler HandlerFunc) {
// 	e.addRoute("GET", pattern, handler)
// }


// func (e *Engine) POST(pattern string, handler HandlerFunc) {
// 	e.addRoute("POST", pattern, handler)
// }


// func (e *Engine) PATCH(pattern string, handler HandlerFunc) {
// 	e.addRoute("PATCH", pattern, handler)
// }


// func (e *Engine) DELETE(pattern string, handler HandlerFunc) {
// 	e.addRoute("DELETE", pattern, handler)
// }

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc

	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := newContext(w, req)
	c.handlers = middlewares
	e.router.handle(c)
}
