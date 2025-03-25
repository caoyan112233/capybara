package capybara

import (
	"encoding/json"
	"net/http"
	"sync"

	"golang.org/x/crypto/acme/autocert"
)

// MIME（Multipurpose Internet Mail Extensions，多用途互联网邮件扩展）是一种用于描述文件类型和格式的标准，
// 最初设计是为了扩展电子邮件的功能，使其能够支持非ASCII文本、二进制文件、多媒体内容等。如今，MIME类型广泛应用于
// 互联网协议（如HTTP）、文件传输、Web开发等领域，是浏览器和服务器之间识别内容类型的核心机制。
const (
	// application type
	APPLICATION_JSON = "application/json" //JSON数据
	APPLICATION_XML  = "application/xml"
	// text type
	TEXT_XML   = "text/xml"
	TEXT_HTML  = "text/html" //HTML网页文件
	TEXT_PLAIN = "text/plain"
)

// 头部常量
const (
	CONTENT_TYPE = "Content-Type"
	// 标识客户端类型，常用于爬虫检测或统计客户端版本。
	USERAGENT = "UserAgent"
	// 传递认证信息（如Bearer Token），常用于API鉴权。
	AUTHORIZATION = "Authorization"
	// 表示请求或响应体的字节长度，如r.Header.Get("Content-Length")
	CONTENT_LENGTH = "ContentLength"
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
	currNode, params := c.router.tree.FindRoute(r.URL.Path)
	if currNode != nil {
		if currNode.method != r.Method {
			sendError(500, w, "Error method")
			return
		}
		// 从池中取出一个context对象
		currContext := c.pool.Get().(*context)
		// 确保方法结束时关闭这个池
		defer c.pool.Put(currContext)
		currContext.Reset()
		currContext.ApplyContext(c, params, w, r)
		currContext.path = currNode.fullPath
		currContext.handler = currNode.handler

		currNode.handler(currContext)
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
