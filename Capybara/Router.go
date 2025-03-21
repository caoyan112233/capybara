package capybara

// **** Router
type Router struct {
	tree        *node
	c           *capybara
	prefix      string
	middlewares []Middlewares
}

func NewRouter() *Router {
	return &Router{
		tree: &node{
			children:  make(map[string]*node),
			wildChild: new(node),
			handler:   make(map[string]HandlerFunc),
		},
	}
}
func (r *Router) GET(path string, handler HandlerFunc, middlewares ...Middlewares) {
	fullPath := joinPath(r.prefix, path)
	if len(r.middlewares) != 0 {
		handler = applyMiddlewares(handler, r.middlewares...)
	}
	r.c.GET(fullPath, handler, middlewares...)
}

func (r *Router) POST(path string, handler HandlerFunc, middlewares ...Middlewares) {
	fullPath := joinPath(r.prefix, path)
	if len(r.middlewares) != 0 {
		handler = applyMiddlewares(handler, r.middlewares...)
	}
	r.c.POST(fullPath, handler, middlewares...)
}

func (r *Router) Use(middlewares ...Middlewares) *Router {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}
