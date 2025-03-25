package capybara

import (
	"encoding/json"
	"net/http"
	"sync"

	"golang.org/x/crypto/acme/autocert"
)

// MIME
const (
	// application type
	APPLICATION_JSON = "application/json"
	APPLICATION_XML  = "application/xml"
	// text type
	TEXT_XML   = "text/xml"
	TEXT_HTML  = "text/html"
	TEXT_PLAIN = "text/plain"
)

const (
	CONTENT_TYPE = "Content-Type"
)

type HandlerFunc func(Context)

// Middlewares
type Middlewares func(HandlerFunc) HandlerFunc

type capybara struct {
	router     *Router
	pool       sync.Pool
	logger     *CapybaraLogger
	TLSManager autocert.Manager
}

// 启动一个capybara实例
func CreateCapybaraInstance() *capybara {
	c := &capybara{
		router: NewRouter(),
		pool: sync.Pool{
			New: func() interface{} {
				// 当池中无可用对象时，自动调用此函数创建新对象
				return new(context)
			}},
		logger: InitLogger(),
		TLSManager: autocert.Manager{
			Prompt: autocert.AcceptTOS,
		},
	}
	c.router.c = c
	return c
}

// 启动非https 的服务
func (c *capybara) Run(addr string) error {
	c.logger.Info(addr + " running")
	err := http.ListenAndServe(addr, c)
	return err
}

// 启动https 的服务
func (c *capybara) RunTLS(addr string, certFile string, keyFile string) error {
	c.logger.Info(addr + " running TLS")
	err := http.ListenAndServeTLS(addr, certFile, keyFile, c)
	return err
}

func (c *capybara) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, params, fields := c.router.tree.FindRoute(r.URL.Path)
	if handler != nil {
		// 从池中取出一个context对象
		currContext := c.pool.Get().(*context)
		// 确保方法结束时关闭这个池
		defer c.pool.Put(currContext)
		currContext.Reset()
		currContext.ApplyContext(c, params, w, r)
		currContext.path = fields[1]
		currContext.handler = handler

		if fields[0] != r.Method {
			sendError(500, w, "Error method")
			return
		}
		c.logger.Info("Call a " + r.Method + " - 200")
		handler(currContext)
	} else {
		sendError(500, w, "Error url")
	}
}

func sendError(code int, w http.ResponseWriter, data interface{}) {
	w.Header().Set(CONTENT_TYPE, APPLICATION_JSON)
	w.WriteHeader(code)
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(map[string]interface{}{"error": data})
}

func (c *capybara) GET(path string, handler HandlerFunc, middlewares ...Middlewares) {
	h := applyMiddlewares(handler, middlewares...)
	c.router.tree.insertRoute(path, "GET", h)
}

func (c *capybara) POST(path string, handler HandlerFunc, middlewares ...Middlewares) {
	h := applyMiddlewares(handler, middlewares...)
	c.router.tree.insertRoute(path, "POST", h)
}

func (c *capybara) DELETE(path string, handler HandlerFunc, middlewares ...Middlewares) {
	h := applyMiddlewares(handler, middlewares...)
	c.router.tree.insertRoute(path, "DELETE", h)
}

func (c *capybara) PUT(path string, handler HandlerFunc, middlewares ...Middlewares) {
	h := applyMiddlewares(handler, middlewares...)
	c.router.tree.insertRoute(path, "PUT", h)
}

func (c *capybara) PATCH(path string, handler HandlerFunc, middlewares ...Middlewares) {
	h := applyMiddlewares(handler, middlewares...)
	c.router.tree.insertRoute(path, "PATCH", h)
}

func (c *capybara) HEAD(path string, handler HandlerFunc, middlewares ...Middlewares) {
	h := applyMiddlewares(handler, middlewares...)
	c.router.tree.insertRoute(path, "HEAD", h)
}

func (c *capybara) OPTIONS(path string, handler HandlerFunc, middlewares ...Middlewares) {
	h := applyMiddlewares(handler, middlewares...)
	c.router.tree.insertRoute(path, "OPTIONS", h)
}

func (c *capybara) TRACE(path string, handler HandlerFunc, middlewares ...Middlewares) {
	h := applyMiddlewares(handler, middlewares...)
	c.router.tree.insertRoute(path, "TRACE", h)
}

func applyMiddlewares(handler HandlerFunc, middlewares ...Middlewares) HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func (c *capybara) Group(prefix string) *Router {
	return &Router{
		prefix:      prefix,
		c:           c,
		tree:        InitNode(),
		middlewares: make([]Middlewares, 0),
	}
}

func Recovery() Middlewares {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx Context) {
			defer func() {
				if err := recover(); err != nil {
					ctx.JSON(502, map[string]string{"error": "StatusInternalServerError"})
				}
			}()
			next(ctx)
		}
	}
}
