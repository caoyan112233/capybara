package capybara

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type HandlerFunc func(Context)

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

func (c *capybara) GET(path string, handler HandlerFunc) {
	c.router.tree.addRoute(path, "GET", handler)
}

func (c *capybara) POST(path string, handler HandlerFunc, middlewares ...Middlewares) {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	c.router.tree.addRoute(path, "POST", handler)
}

// **** context
type Context interface {
	JSON(data interface{})
	Request() *http.Request
	GetHeader(key string) string
	Set(key string, value interface{})
	Get(key string) interface{}
	Bind(data interface{}) (err error)
}
type context struct {
	w    http.ResponseWriter
	r    *http.Request
	data map[string]interface{}
}

func (c *context) JSON(data interface{}) {
	jsonEncoder := json.NewEncoder(c.w)
	jsonEncoder.Encode(data)
}
func (c *context) Request() *http.Request {
	return c.r
}
func (c *context) GetHeader(key string) string {

	return c.r.Header.Get(key)
}

func (c *context) Get(key string) interface{} {
	return c.data[key]
}

func (c *context) Set(key string, value interface{}) {
	c.data[key] = value
}

func (c *context) Bind(data interface{}) (err error) {
	// 读取请求体
	body, err := io.ReadAll(c.r.Body)
	if err != nil {
		return err
	}
	defer c.r.Body.Close()

	// 如果请求体i为空，返回错误
	if len(body) == 0 {
		return errors.New("request body is empty")
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}
	return err
}

// **** Router
type Router struct {
	tree        *node
	c           *capybara
	prefix      string
	middlewares []Middlewares
}

func (r *Router) GET(path string, handler HandlerFunc) {
	fullPath := joinPath(r.prefix, path)
	wrappedHanler := applyMiddlewares(handler, r.middlewares)
	r.c.GET(fullPath, wrappedHanler)
}

func (r *Router) POST(path string, handler HandlerFunc, middlewares ...Middlewares) {
	fullPath := joinPath(r.prefix, path)
	wrappedHanler := applyMiddlewares(handler, r.middlewares)
	r.c.POST(fullPath, wrappedHanler, middlewares...)
}

func applyMiddlewares(handler HandlerFunc, middlewares []Middlewares) HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

func (r *Router) Use(middlewares ...Middlewares) *Router {
	r.middlewares = append(r.middlewares, middlewares...)
	return r
}

// 辅助函数：合并前缀和路径，处理斜杠问题
func joinPath(prefix, path string) string {
	if prefix == "" {
		return path
	}
	if path == "" {
		return prefix
	}
	return strings.TrimSuffix(prefix, "/") + "/" + strings.TrimPrefix(path, "/")
}

func (c *capybara) Group(prefix string) *Router {
	return &Router{
		prefix: prefix,
		c:      c,
	}
}

// **** route tree
type node struct {
	children  map[string]*node
	wildChild *node
	isWild    bool
	paramKey  string
	handler   map[string]HandlerFunc
}

func splitPath(path string) []string {
	if path == "" {
		return []string{}
	}
	return strings.Split(strings.TrimPrefix(path, "/"), "/")
}

// 增加一条路由
func (n *node) addRoute(path string, method string, handler HandlerFunc) {
	segs := splitPath(path)
	for _, seg := range segs {
		if seg[0] == ':' || seg[0] == '*' {
			n.wildChild = &node{
				children:  make(map[string]*node),
				wildChild: new(node),
				isWild:    true,
				paramKey:  seg[1:],
				handler:   make(map[string]HandlerFunc)}
			n = n.wildChild
		} else {
			if n.children == nil {
				n.children = make(map[string]*node)
			}
			if _, ok := n.children[seg]; !ok {
				n.children[seg] = &node{
					children:  make(map[string]*node),
					wildChild: new(node),
					isWild:    false,
					paramKey:  seg,
					handler:   make(map[string]HandlerFunc)}
			}
			n = n.children[seg]
		}
	}
	n.handler[method] = handler
}

// 匹配一条路由
func (n *node) matchRoute(path string) *node {
	segs := splitPath(path)
	current := n
	for _, seg := range segs {
		if child, ok := current.children[seg]; ok {
			current = child
			continue
		}
		if current.wildChild != nil {
			wildChild := current.wildChild
			if wildChild.isWild {
				current = wildChild
				break
			} else {
				current = wildChild
				continue
			}
		}
		return nil
	}
	return current
}

// Middlewares
type Middlewares func(HandlerFunc) HandlerFunc
