package capybara

// **** Router
type Router struct {
	c           *capybara
	tree        *node
	prefix      string
	middlewares []Middlewares
}

func NewRouter() *Router {
	return &Router{
		tree:        InitNode(),
		c:           nil,
		prefix:      "",
		middlewares: make([]Middlewares, 0),
	}
}

// http 请求组的  GET 方法
func (r *Router) GET(path string, handler HandlerFunc, middlewares ...Middlewares) {
	fullPath := joinPath(r.prefix, path)
	if len(r.middlewares) != 0 {
		handler = applyMiddlewares(handler, r.middlewares...)
	}
	r.c.GET(fullPath, handler, middlewares...)
}

// http 请求组的  POST 方法
func (r *Router) POST(path string, handler HandlerFunc, middlewares ...Middlewares) {
	fullPath := joinPath(r.prefix, path)
	if len(r.middlewares) != 0 {
		handler = applyMiddlewares(handler, r.middlewares...)
	}
	r.c.POST(fullPath, handler, middlewares...)
}

// http 请求组的  DELETE 方法
func (r *Router) DELETE(path string, handler HandlerFunc, middlewares ...Middlewares) {
	fullPath := joinPath(r.prefix, path)
	if len(r.middlewares) != 0 {
		handler = applyMiddlewares(handler, r.middlewares...)
	}
	r.c.POST(fullPath, handler, middlewares...)
}

// http 请求组的  HEAD 方法
func (r *Router) HEAD(path string, handler HandlerFunc, middlewares ...Middlewares) {
	fullPath := joinPath(r.prefix, path)
	if len(r.middlewares) != 0 {
		handler = applyMiddlewares(handler, r.middlewares...)
	}
	r.c.POST(fullPath, handler, middlewares...)
}

// http 请求组的  OPTIONS 方法
func (r *Router) OPTIONS(path string, handler HandlerFunc, middlewares ...Middlewares) {
	fullPath := joinPath(r.prefix, path)
	if len(r.middlewares) != 0 {
		handler = applyMiddlewares(handler, r.middlewares...)
	}
	r.c.POST(fullPath, handler, middlewares...)
}

// http 请求组的  PATCH 方法
func (r *Router) PATCH(path string, handler HandlerFunc, middlewares ...Middlewares) {
	fullPath := joinPath(r.prefix, path)
	if len(r.middlewares) != 0 {
		handler = applyMiddlewares(handler, r.middlewares...)
	}
	r.c.POST(fullPath, handler, middlewares...)
}

// http 请求组的  PUT 方法
func (r *Router) PUT(path string, handler HandlerFunc, middlewares ...Middlewares) {
	fullPath := joinPath(r.prefix, path)
	if len(r.middlewares) != 0 {
		handler = applyMiddlewares(handler, r.middlewares...)
	}
	r.c.POST(fullPath, handler, middlewares...)
}

// http 请求组的  TRACE 方法
func (r *Router) TRACE(path string, handler HandlerFunc, middlewares ...Middlewares) {
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
