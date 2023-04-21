package pulse

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type handlerFunc func(ctx *Context) error

type Context struct {
	ResponseWriter http.ResponseWriter
	Response       http.Response
	Request        *http.Request
	Params         map[string]string
	paramValues    []string
	handlers       []handlerFunc
	handlerIdx     int
	Cookies        []*http.Cookie
}

func (c *Context) Write(p []byte) (n int, err error) {
	return c.ResponseWriter.Write(p)
}

type Cookie struct {
	Name     string
	Value    string
	Path     string
	Domain   string
	MaxAge   int
	Expires  time.Time
	Secure   bool
	HTTPOnly bool
	SameSite http.SameSite
}

// NewContext returns a new Context.
func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		ResponseWriter: w,
		Request:        req,
		Params:         make(map[string]string),
		paramValues:    make([]string, 0, 10),
		handlers:       nil,
		handlerIdx:     -1,
	}
}

// WithParams sets the params for the context.
func (c *Context) WithParams(params map[string]string) *Context {
	c.Params = params
	return c
}

// Param returns the param value for the given key.
func (c *Context) Param(key string) string {
	return c.Params[key]
}

// Query returns the query value for the given key.
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// String sets the response body to the given string.
func (c *Context) String(value string) {
	_, err := c.ResponseWriter.Write([]byte(value))
	if err != nil {
		return
	}
}

// SetValue create a middleware that adds a value to the context
func (c *Context) SetValue(key interface{}, value interface{}) {
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), key, value))
}

// GetValue returns the value for the given key.
func (c *Context) GetValue(key string) string {
	return c.Request.Context().Value(key).(string)
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
	http.SetCookie(c.ResponseWriter, &http.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		Expires:  cookie.Expires,
		MaxAge:   cookie.MaxAge,
		Secure:   cookie.Secure,
		HttpOnly: cookie.HTTPOnly,
		SameSite: cookie.SameSite,
	})

	c.Cookies = append(c.Cookies, &http.Cookie{
		Name:     cookie.Name,
		Value:    cookie.Value,
		Path:     cookie.Path,
		Domain:   cookie.Domain,
		Expires:  cookie.Expires,
		MaxAge:   cookie.MaxAge,
		Secure:   cookie.Secure,
		HttpOnly: cookie.HTTPOnly,
		SameSite: cookie.SameSite,
	})
}

// GetCookie returns the value of the cookie with the given name.
func (c *Context) GetCookie(name string) string {
	for _, cookie := range c.Cookies {
		if cookie.Name == name {
			return cookie.Value
		}
	}
	return ""
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

// SetResponseHeader sets the http header value to the given key.
func (c *Context) SetResponseHeader(key, value string) {
	c.ResponseWriter.Header().Set(key, value)
}

// GetResponseHeader returns the http header value for the given key.
func (c *Context) GetResponseHeader(key string) string {
	return c.ResponseWriter.Header().Get(key)
}

// SetRequestHeader SetResponseHeader sets the http header value to the given key.
func (c *Context) SetRequestHeader(key, value string) {
	c.Request.Header.Set(key, value)
}

// GetRequestHeader GetResponseHeader returns the http header value for the given key.
func (c *Context) GetRequestHeader(key string) string {
	return c.Request.Header.Get(key)
}

// SetContentType sets the Content-Type header in the response to the given value.
func (c *Context) SetContentType(value string) {
	c.ResponseWriter.Header().Set("Content-Type", value)
}

// Accepts checks if the specified content types are acceptable.
func (c *Context) Accepts(types ...string) string {
	acceptHeader := c.GetRequestHeader("Accept")
	if acceptHeader == "" {
		return ""
	}

	acceptedMediaTypes := strings.Split(acceptHeader, ",")

	for _, t := range types {
		for _, a := range acceptedMediaTypes {
			a = strings.TrimSpace(a)
			if strings.HasPrefix(a, t+"/") || a == "*/*" || a == t {
				return t
			}
		}
	}

	return ""
}

// Status sets the response status code.
func (c *Context) Status(code int) {
	c.ResponseWriter.WriteHeader(code)
}

// JSON sets the response body to the given JSON representation.
func (c *Context) JSON(code int, obj interface{}) ([]byte, error) {
	c.ResponseWriter.Header().Set("Content-Type", "application/json")
	c.Status(code)
	jsonBody, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	if _, err := c.ResponseWriter.Write(jsonBody); err != nil {
		return nil, err
	}

	return jsonBody, nil
}

func (c *Context) BodyParser(v interface{}) error {
	return json.NewDecoder(c.Request.Body).Decode(v)
}
