package pulse

type MiddlewareFunc func(handler Handler) Handler

type Middleware interface {
	Middleware(handler Handler) Handler
	Handle(ctx *Context, next Handler) error
}

func (m MiddlewareFunc) Middleware(handler Handler) Handler {
	return m(handler)
}

func (r *Router) Use(method string, middlewares ...interface{}) {
	for _, middleware := range middlewares {
		if middlewareFunc, ok := middleware.(MiddlewareFunc); ok {
			r.middlewares[method] = append(r.middlewares[method], middlewareFunc)
		} else if middleware, ok := middleware.(Middleware); ok {
			r.middlewares[method] = append(r.middlewares[method], middleware)
		}
	}
}

func CORSMiddleware() MiddlewareFunc {
	return func(handler Handler) Handler {
		return func(ctx *Context) error {
			ctx.SetResponseHeader("Access-Control-Allow-Origin", "*")
			ctx.SetResponseHeader("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			ctx.SetResponseHeader("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			return handler(ctx)
		}
	}
}

func (m MiddlewareFunc) Handle(ctx *Context, next Handler) error {
	h := m(next)
	return h(ctx)
}
