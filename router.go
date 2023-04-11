package pulse

import (
	"fmt"
	"github.com/gopulse/pulse/constants"
	"github.com/valyala/fasthttp"
	"strings"
	"time"
)

type Handler func(ctx *Context) error

type Route struct {
	Path       string
	ParamNames []string
	Handlers   []Handler
}

type Router struct {
	routes          map[string][]*Route
	notFoundHandler Handler
	middlewares     map[string][]Middleware
}

type Static struct {
	Root          string
	Compress      bool
	ByteRange     bool
	IndexName     string
	CacheDuration time.Duration
}

func NewRouter() *Router {
	return &Router{
		routes:      make(map[string][]*Route),
		middlewares: make(map[string][]Middleware),
	}
}

func (r *Router) add(method, path string, handlers []Handler) {
	route := &Route{
		Path:     path,
		Handlers: handlers,
	}

	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, constants.ParamSign) {
			route.ParamNames = append(route.ParamNames, strings.TrimPrefix(part, constants.ParamSign))
		}
	}
	route.Path = strings.Join(parts, "/")

	r.routes[method] = append(r.routes[method], route)
}

func (r *Router) find(method, path string) []Handler {
	routes, ok := r.routes[method]
	if !ok {
		return nil
	}

	for _, route := range routes {
		if matches, params := route.match(path); matches {
			c := NewContext(nil, nil)
			c.params = params
			return r.applyMiddleware(route.Handlers, method)
		}
	}

	return nil
}

func (r *Router) applyMiddleware(handlers []Handler, method string) []Handler {
	for i := len(r.middlewares[method]) - 1; i >= 0; i-- {
		middleware := r.middlewares[method][i]
		for j := len(handlers) - 1; j >= 0; j-- {
			handler := handlers[j]
			handlers[j] = func(ctx *Context) error {
				return middleware.Handle(ctx, handler)
			}
		}
	}
	return handlers
}

func RouterHandler(router *Router) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())
		method := string(ctx.Method())
		handlers := router.find(method, path)
		if handlers == nil {
			ctx.Error("Page not found", fasthttp.StatusNotFound)
			return
		}

		params := make(map[string]string)
		for _, route := range router.routes[method] {
			if matches, params := route.match(path); matches {
				c := NewContext(ctx, nil).WithParams(params)
				for _, h := range router.applyMiddleware(route.Handlers, method) {
					err := h(c)
					if err != nil {
						ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
						return
					}
				}
				return
			}
		}

		c := NewContext(ctx, params)

		for _, h := range router.applyMiddleware(handlers, method) {
			err := h(c)
			if err != nil {
				ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
				return
			}
		}
	}
}

func (r *Route) match(path string) (bool, map[string]string) {
	parts := strings.Split(path, "/")
	routeParts := strings.Split(r.Path, "/")

	if len(parts) != len(routeParts) {
		return false, nil
	}

	params := make(map[string]string)
	for i, part := range routeParts {
		if strings.HasPrefix(part, constants.ParamSign) {
			paramName := strings.TrimPrefix(part, constants.ParamSign)
			params[paramName] = parts[i]
		} else if part == constants.WildcardSign {
			return true, params
		} else if part != parts[i] {
			return false, nil
		}
	}

	return true, params
}

func (options *Static) notFoundHandler(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.SetContentType("text/plain; charset=utf-8")
	_, err := fmt.Fprintf(ctx, "404 Not Found")
	if err != nil {
		return
	}
}

func (options *Static) pathRewrite(ctx *fasthttp.RequestCtx) []byte {
	path := ctx.Path()

	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	// Remove the last part of the path
	parts := strings.Split(string(path), "/")
	if len(parts) > 1 {
		parts = parts[:len(parts)-1]
	}
	path = []byte(strings.Join(parts, "/"))

	if options.IndexName != "" {
		// Append the index file name to the path
		path = append(path, '/')
		path = append(path, options.IndexName...)
	}

	return path
}

func (r *Router) Static(prefix, root string, options *Static) {
	if options == nil {
		options = &Static{}
	}
	if options.Root == "" {
		options.Root = root
	}
	fs := fasthttp.FS{
		Root:               options.Root,
		IndexNames:         []string{options.IndexName},
		PathRewrite:        options.pathRewrite,
		GenerateIndexPages: false,
		Compress:           options.Compress,
		AcceptByteRange:    options.ByteRange,
		CacheDuration:      options.CacheDuration,
		PathNotFound:       options.notFoundHandler, // Set custom error handler for undefined routes
	}

	r.Get(prefix, func(c *Context) error {
		fsHandler := fs.NewRequestHandler()
		fsHandler(c.RequestCtx)
		return nil
	})
}
