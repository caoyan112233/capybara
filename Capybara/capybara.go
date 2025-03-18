package capybara

import (
	"encoding/json"
	"net/http"
)

type HandlerFunc func(Context)

// Middlewares
type Middlewares func(HandlerFunc) HandlerFunc

type capybara struct {
	router *Router
}

func New() *capybara {
	c := &capybara{
		router: &Router{
			tree: &node{
				children:  make(map[string]*node),
				wildChild: new(node),
				isWild:    false,
				handler:   make(map[string]HandlerFunc),
			}}}
	return c
}

func (c *capybara) Run(addr string) {
	http.ListenAndServe(addr, c)
}

func (c *capybara) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	currentNode := c.router.tree.matchRoute(r.URL.Path)
	if currentNode != nil {
		currContext := &context{
			w:    w,
			r:    r,
			data: make(map[string]interface{}),
		}
		if currentNode.handler[r.Method] != nil {
			currentNode.handler[r.Method](currContext)
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
