package capybara

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
)

// **** context
type Context interface {
	JSON(code int, data interface{}) error
	XML(code int, data interface{}) error
	String(code int, s string) error
	HTML(code int, html string) error

	Request() *http.Request
	GetHeader(key string) string
	Set(key string, value interface{})
	Get(key string) interface{}
	Bind(data interface{}) (err error)

	Param(name string) string
	Path() string
	Handler() HandlerFunc
	Cookie(name string) (*http.Cookie, error)
	Cookies() []*http.Cookie
}

type context struct {
	w       http.ResponseWriter
	r       *http.Request
	data    map[string]interface{}
	capa    *capybara
	params  map[string]string
	path    string
	handler HandlerFunc
}

// 应用到当前的 context
func (c *context) ApplyContext(cap *capybara, params map[string]string, w http.ResponseWriter, r *http.Request) {
	c.capa = cap
	c.w = w
	c.r = r
	c.data = make(map[string]interface{})
	c.params = params
}

// 发送JSON格式的文件
func (c *context) JSON(code int, data interface{}) error {
	jsonEncoder := json.NewEncoder(c.w)
	err := jsonEncoder.Encode(data)
	if err != nil {
		return c.String(http.StatusInternalServerError, "解析json出错")
	}
	return nil
}

// 发送String格式的文件
func (c *context) String(code int, s string) error {
	c.w.Header().Set("Content-Type", "text/plain")
	c.w.WriteHeader(code)
	_, err := c.w.Write([]byte(s))
	return err
}

// 发送 XML格式的文件
func (c *context) XML(code int, data interface{}) error {
	c.w.Header().Set("Content-Type", "application/xml")
	c.w.WriteHeader(code)
	xmlEncoder := xml.NewEncoder(c.w)
	err := xmlEncoder.Encode(data)
	if err != nil {
		c.String(http.StatusInternalServerError, "解析xml出错")
	}
	return nil
}

// 发送HTNLi格式的文件
func (c *context) HTML(code int, html string) error {
	c.w.Header().Set("Content-Type", "text/html")
	c.w.WriteHeader(code)
	_, err := c.w.Write([]byte(html))
	return err
}

// 获取一个 路由中的某个指定的参数
//
//	例如：
//
// http://localhost:8080/user/123/post/456
//
//	/user/:id/post/:post_id
//
// id ： 123
// post_id : 456
func (c *context) Param(name string) string {

	if _, ok := c.params[name]; !ok {
		return ""
	}
	return c.params[name]
}

// 获取路由路径
//
// /user/:id/post/:post_id
func (c *context) Path() string {
	return c.path
}

// 获取路由函数
func (c *context) Handler() HandlerFunc {
	return c.handler
}

// 返回Cookie
func (c *context) Cookie(name string) (*http.Cookie, error) {
	return c.r.Cookie(name)
}

func (c *context) Cookies() []*http.Cookie {
	return c.r.Cookies()
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

	// 如果请求体为空，返回错误
	if len(body) == 0 {
		return errors.New("request body is empty")
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return err
	}
	return err
}
