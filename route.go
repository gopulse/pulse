package routing

import "github.com/gopulse/pulse-router/constants"

type route struct {
	method  string
	path    string
	handler Handler
}

type RouteStore interface {
	Add(key string, data interface{}) int
	Get(path string, pvalues []string) (data interface{}, pnames []string)
}

func newRoute(method, path string, handler Handler) *route {
	return &route{method, path, handler}
}

// Name sets the name of the route
func (r *route) Name() string {
	return r.method + r.path
}

// Get adds the route to the router with the GET method
func (r *Router) Get(path string, handlers ...Handler) {
	r.add(constants.GetMethod, path, handlers)
}

// Post adds the route to the router with the POST method
func (r *Router) Post(path string, handlers ...Handler) {
	r.add(constants.PostMethod, path, handlers)
}

// Put adds the route to the router with the PUT method
func (r *Router) Put(path string, handlers ...Handler) {
	r.add(constants.PutMethod, path, handlers)
}
