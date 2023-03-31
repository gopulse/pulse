package routing

import "github.com/valyala/fasthttp"

type handlerFunc func(ctx *context) error

type context struct {
	ctx         *fasthttp.RequestCtx
	params      map[string]string
	paramValues []string
	handlers    []handlerFunc
	handlerIdx  int
}

// Context returns the fasthttp.RequestCtx
func (c *context) Context() *fasthttp.RequestCtx {
	return c.ctx
}

// Param returns the param value for the given key.
func (c *context) Param(key string) string {
	return c.params[key]
}

// String sets the response body to the given string.
func (c *context) String(value string) error {
	c.ctx.SetBodyString(value)
	return nil
}

// SetData sets the http header value to the given key.
func (c *context) SetData(key string, value interface{}) {
	c.ctx.Response.Header.Set(key, value.(string))
}

// GetData returns the http header value for the given key.
func (c *context) GetData(key string) string {
	return string(c.ctx.Response.Header.Peek(key))
}

// Next calls the next handler in the chain.
func (c *context) Next() error {
	c.handlerIdx++
	if c.handlerIdx < len(c.handlers) {
		return c.handlers[c.handlerIdx](c)
	}
	return nil
}

// Reset resets the context.
func (c *context) Reset() {
	c.handlerIdx = -1
}

// Abort aborts the chain.
func (c *context) Abort() {
	c.handlerIdx = len(c.handlers)
}
