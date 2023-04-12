package pulse

import (
	"bytes"
	"github.com/gopulse/pulse/utils"
	"github.com/valyala/fasthttp"
	"time"
)

type handlerFunc func(ctx *Context) error

type Context struct {
	RequestCtx  *fasthttp.RequestCtx
	params      map[string]string
	paramValues []string
	handlers    []handlerFunc
	handlerIdx  int
	Cookies     *fasthttp.Cookie
}

type Cookie struct {
	Name        string    `json:"name"`
	Value       string    `json:"value"`
	Path        string    `json:"path"`
	Domain      string    `json:"domain"`
	MaxAge      int       `json:"max_age"`
	Expires     time.Time `json:"expires"`
	Secure      bool      `json:"secure"`
	HTTPOnly    bool      `json:"http_only"`
	SameSite    string    `json:"same_site"`
	SessionOnly bool      `json:"session_only"`
}

// NewContext returns a new Context.
func NewContext(ctx *fasthttp.RequestCtx, params map[string]string) *Context {
	return &Context{
		RequestCtx:  ctx,
		params:      params,
		paramValues: make([]string, 0, 10),
		handlers:    nil,
		handlerIdx:  -1,
	}
}

// Context returns the fasthttp.RequestCtx
func (c *Context) Context() *fasthttp.RequestCtx {
	return c.RequestCtx
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

// Query returns the query value for the given key.
func (c *Context) Query(key string) string {
	return string(c.RequestCtx.QueryArgs().Peek(key))
}

// String sets the response body to the given string.
func (c *Context) String(value string) {
	if c.RequestCtx.Response.Body() == nil {
		c.RequestCtx.Response.SetBodyString(value)
	} else {
		buf := bytes.NewBuffer(c.RequestCtx.Response.Body())
		buf.WriteString(value)
		c.RequestCtx.Response.SetBody(buf.Bytes())
	}
}

// SetData sets the http header value to the given key.
func (c *Context) SetData(key string, value interface{}) {
	c.SetResponseHeader(key, value.(string))
}

// GetData returns the http header value for the given key.
func (c *Context) GetData(key string) string {
	return string(c.RequestCtx.Response.Header.Peek(key))
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

// SetCookie sets a cookie with the given name, value, and options.
func (c *Context) SetCookie(cookie *Cookie) {
	acCookie := fasthttp.AcquireCookie()

	acCookie.SetKey(cookie.Name)
	acCookie.SetValue(cookie.Value)
	acCookie.SetPath(cookie.Path)
	acCookie.SetDomain(cookie.Domain)
	acCookie.SetSecure(cookie.Secure)
	acCookie.SetHTTPOnly(cookie.HTTPOnly)
	acCookie.SetSecure(cookie.Secure)
	if !cookie.SessionOnly {
		acCookie.SetMaxAge(cookie.MaxAge)
		acCookie.SetExpire(cookie.Expires)
	}

	switch utils.ToLower(cookie.SameSite) {
	case string(rune(fasthttp.CookieSameSiteStrictMode)):
		acCookie.SetSameSite(fasthttp.CookieSameSiteStrictMode)
	case string(rune(fasthttp.CookieSameSiteNoneMode)):
		acCookie.SetSameSite(fasthttp.CookieSameSiteNoneMode)
	case string(rune(fasthttp.CookieSameSiteDisabled)):
		acCookie.SetSameSite(fasthttp.CookieSameSiteDisabled)
	default:
		acCookie.SetSameSite(fasthttp.CookieSameSiteDefaultMode)
	}

	c.RequestCtx.Response.Header.SetCookie(acCookie)
	fasthttp.ReleaseCookie(acCookie)
}

// GetCookie returns the value of the cookie with the given name.
func (c *Context) GetCookie(name string) string {
	cookie := c.RequestCtx.Request.Header.Cookie(name)
	if cookie == nil {
		return ""
	}
	return string(cookie)
}

// ClearCookie deletes the cookie with the given name.
func (c *Context) ClearCookie(name string) {
	c.SetCookie(&Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   false,
		HTTPOnly: true,
	})
}

// SetHeader sets the http header value to the given key.
func (c *Context) SetHeader(key, value string) {
	c.RequestCtx.Response.Header.Set(key, value)
}

// GetHeader returns the http header value for the given key.
func (c *Context) GetHeader(key string) string {
	return string(c.RequestCtx.Request.Header.Peek(key))
}

// Accepts returns true if the specified type(s) is acceptable, otherwise false.
func (c *Context) Accepts(types ...string) string {
	return ""
}

// Status sets the response status code.
func (c *Context) Status(code int) {
	c.RequestCtx.Response.SetStatusCode(code)
}

// JSON sets the response body to the given JSON representation.
func (c *Context) JSON(code int, obj interface{}) {
	c.RequestCtx.Response.Header.SetContentType("application/json")
	c.RequestCtx.Response.SetStatusCode(code)
	c.RequestCtx.Response.SetBodyString(utils.ToJSON(obj))
}
