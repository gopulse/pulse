package pulse

import (
	"github.com/gopulse/pulse/constants"
	"net/http"
	"strings"
	"time"
)

type Handler func(ctx *Context) error

type Route struct {
	Method     string
	Path       string
	Handlers   []Handler
	ParamNames []string
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

func (r *Router) Add(method, path string, handlers ...Handler) {
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

func (r *Router) Find(method, path string) []Handler {
	routes, ok := r.routes[method]
	if !ok {
		return nil
	}

	for _, route := range routes {
		if matches, params := route.match(path); matches {
			c := NewContext(nil, nil)
			c.Params = params
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

func RouterHandler(router *Router) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		path := req.URL.Path
		method := req.Method
		handlers := router.Find(method, path)

		c := NewContext(w, req)
		for _, h := range handlers {
			err := h(c)
			if err != nil {
				break
			}
		}
	}
}

func (r *Route) match(path string) (bool, map[string]string) {
	parts := strings.Split(path, "/")
	routeParts := strings.Split(r.Path, "/")

	if strings.HasSuffix(path, "/") {
		parts = parts[:len(parts)-1]
	}

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

func (options *Static) notFoundHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := w.Write([]byte("404 Not Found"))
	if err != nil {
		return
	}
}

func (options *Static) PathRewrite(r *http.Request) []byte {
	path := r.URL.Path

	if len(path) > 1 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	// Remove the last part of the path
	parts := strings.Split(path, "/")
	if len(parts) > 1 {
		parts = parts[:len(parts)-1]
	}
	path = strings.Join(parts, "/")

	if options.IndexName != "" {
		// Append the index file name to the path
		path += "/"
		path += options.IndexName
	}

	return []byte(path)
}

func (r *Router) Static(prefix, root string, options *Static) {
	if options == nil {
		options = &Static{}
	}
	if options.Root == "" {
		options.Root = root
	}
	fs := http.FileServer(http.Dir(options.Root))

	handler := http.StripPrefix(prefix, fs)

	r.Get(prefix, func(ctx *Context) error {
		handler.ServeHTTP(ctx.ResponseWriter, ctx.Request)
		return nil
	})
}
