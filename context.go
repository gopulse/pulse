package routing

import "github.com/valyala/fasthttp"

type handlerFunc func(ctx *Context) error

type Context struct {
	ctx         *fasthttp.RequestCtx
	params      map[string]string
	paramValues []string
	handlers    []handlerFunc
	handlerIdx  int
}

// Context returns the fasthttp.RequestCtx
func (c *Context) Context() *fasthttp.RequestCtx {
	return c.ctx
}

// Param returns the param value for the given key.
func (c *Context) Param(key string) string {
	return c.params[key]
}

// String sets the response body to the given string.
func (c *Context) String(value string) error {
	c.ctx.SetBodyString(value)
	return nil
}

// SetData sets the http header value to the given key.
func (c *Context) SetData(key string, value interface{}) {
	c.ctx.Response.Header.Set(key, value.(string))
}

// GetData returns the http header value for the given key.
func (c *Context) GetData(key string) string {
	return string(c.ctx.Response.Header.Peek(key))
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
