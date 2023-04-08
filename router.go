package routing

import (
	"github.com/valyala/fasthttp"
	"strings"
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

func New() *Router {
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
		if strings.HasPrefix(part, ":") {
			route.ParamNames = append(route.ParamNames, strings.TrimPrefix(part, ":"))
		}
	}
	route.Path = strings.Join(parts, "/")

	r.routes[method] = append(r.routes[method], route)
}

func (r *Router) find(method, path string) []Handler {
	routes := r.routes[method]
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
			ctx.Error("Not found", fasthttp.StatusNotFound)
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
		if strings.HasPrefix(part, ":") {
			paramName := strings.TrimPrefix(part, ":")
			params[paramName] = parts[i]
		} else if part != parts[i] {
			return false, nil
		}
	}

	return true, params
}
