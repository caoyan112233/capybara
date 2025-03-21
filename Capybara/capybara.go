package capybara

import (
	"encoding/json"
	"net/http"
	"sync"
)

type HandlerFunc func(Context)

// Middlewares
type Middlewares func(HandlerFunc) HandlerFunc

type capybara struct {
	router *Router
	pool   sync.Pool
}

func New() *capybara {
	c := &capybara{
		router: NewRouter(),
		pool: sync.Pool{
			New: func() interface{} {
				// 当池中无可用对象时，自动调用此函数创建新对象
				return new(context) // 可以是结构体、切片等任意类型
			}},
	}
	return c
}

func (c *capybara) Run(addr string) {
	http.ListenAndServe(addr, c)
}

func (c *capybara) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	currentNode := c.router.tree.matchRoute(r.URL.Path)
	if currentNode != nil {
		// 从h池中取出一个context对象
		currContext := c.pool.Get().(*context)
		currContext.capa = c
		currContext.w = w
		currContext.r = r
		currContext.data = make(map[string]interface{})

		if currentNode.handler[r.Method] != nil {
			currentNode.handler[r.Method](currContext)
			c.pool.Put(currContext)
		} else {
			sendError(w, "Error method")
		}
	} else {
		sendError(w, "Error url")
	}
}

func sendError(w http.ResponseWriter, data interface{}) {
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(data)
}

func (c *capybara) GET(path string, handler HandlerFunc, middlewares ...Middlewares) {
	h := applyMiddlewares(handler, middlewares...)
	c.router.tree.addRoute(path, "GET", h)
}

func (c *capybara) POST(path string, handler HandlerFunc, middlewares ...Middlewares) {
	h := applyMiddlewares(handler, middlewares...)
	c.router.tree.addRoute(path, "POST", h)
}

func applyMiddlewares(handler HandlerFunc, middlewares ...Middlewares) HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func (c *capybara) Group(prefix string) *Router {
	return &Router{
		prefix: prefix,
		c:      c,
	}
}
