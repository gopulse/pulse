package routing

type Handler func(ctx *Context) error

type Router struct {
	routes          map[string]*route
	stores          map[string]RouteStore
	notFoundHandler Handler
	maxParams       int
}

func CombineHandlers(handlers ...Handler) Handler {
	return func(ctx *Context) error {
		for _, handler := range handlers {
			if err := handler(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}

func New() *Router {
	return &Router{
		routes: make(map[string]*route),
		stores: make(map[string]RouteStore),
	}
}

func (r *Router) add(method, path string, handlers []Handler) {
	s := r.stores[method]
	if s == nil {
		s = newStore()
		r.stores[method] = s
	}
	if n := s.Add(path, handlers); n > r.maxParams {
		r.maxParams = n
	}
	r.routes[method+path] = &route{method, path, CombineHandlers(handlers...)}
}

func (r *Router) find(method, path string, pvalues []string) (handlers []Handler, pnames []string) {
	var hh interface{}
	if store := r.stores[method]; store != nil {
		hh, pnames = store.Get(path, pvalues)
	}
	if hh != nil {
		return hh.([]Handler), pnames
	}
	return nil, nil
}
