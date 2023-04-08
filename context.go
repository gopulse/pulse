package routing

import (
	"bytes"
	"github.com/valyala/fasthttp"
)

type handlerFunc func(ctx *Context) error

type Context struct {
	Ctx         *fasthttp.RequestCtx
	params      map[string]string
	paramValues []string
	handlers    []handlerFunc
	handlerIdx  int
}

// NewContext returns a new Context.
func NewContext(ctx *fasthttp.RequestCtx, params map[string]string) *Context {
	return &Context{
		Ctx:         ctx,
		params:      params,
		paramValues: make([]string, 0, 10),
		handlers:    nil,
		handlerIdx:  -1,
	}
}

// Context returns the fasthttp.RequestCtx
func (c *Context) Context() *fasthttp.RequestCtx {
	return c.Ctx
}

// WithParams sets the params for the context.
func (c *Context) WithParams(params map[string]string) *Context {
	c.params = params
	return c
}

// Param returns the param value for the given key.
func (c *Context) Param(key string) string {
	return c.params[key]
}

// String sets the response body to the given string.
func (c *Context) String(value string) {
	if c.Ctx.Response.Body() == nil {
		c.Ctx.Response.SetBodyString(value)
	} else {
		buf := bytes.NewBuffer(c.Ctx.Response.Body())
		buf.WriteString(value)
		c.Ctx.Response.SetBody(buf.Bytes())
	}
}

// SetData sets the http header value to the given key.
func (c *Context) SetData(key string, value interface{}) {
	c.Ctx.Response.Header.Set(key, value.(string))
}

// GetData returns the http header value for the given key.
func (c *Context) GetData(key string) string {
	return string(c.Ctx.Response.Header.Peek(key))
}

// Next calls the next handler in the chain.
func (c *Context) Next() error {
	c.handlerIdx++
	if c.handlerIdx < len(c.handlers) {
		return c.handlers[c.handlerIdx](c)
	}
	return nil
}

// Reset resets the Context.
func (c *Context) Reset() {
	c.handlerIdx = -1
}

// Abort aborts the chain.
func (c *Context) Abort() {
	c.handlerIdx = len(c.handlers)
}
