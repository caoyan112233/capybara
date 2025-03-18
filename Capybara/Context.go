package capybara

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

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
