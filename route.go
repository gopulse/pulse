package routing

import "github.com/gopulse/pulse-router/constants"

type route struct {
	method  string
	path    string
	handler Handler
}

type RouteStore interface {
	Add(key string, data interface{}) int
	Get(path string) (data interface{})
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

// Delete adds the route to the router with the DELETE method
func (r *Router) Delete(path string, handlers ...Handler) {
	r.add(constants.DeleteMethod, path, handlers)
}

// Patch adds the route to the router with the PATCH method
func (r *Router) Patch(path string, handlers ...Handler) {
	r.add(constants.PatchMethod, path, handlers)
}

// Head adds the route to the router with the HEAD method
func (r *Router) Head(path string, handlers ...Handler) {
	r.add(constants.HeadMethod, path, handlers)
}

// Options adds the route to the router with the OPTIONS method
func (r *Router) Options(path string, handlers ...Handler) {
	r.add(constants.OptionsMethod, path, handlers)
}

// Connect adds the route to the router with the CONNECT method
func (r *Router) Connect(path string, handlers ...Handler) {
	r.add(constants.ConnectMethod, path, handlers)
}

// Trace adds the route to the router with the TRACE method
func (r *Router) Trace(path string, handlers ...Handler) {
	r.add(constants.TraceMethod, path, handlers)
}
