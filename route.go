package pulse

import (
	"net/http"
)

type Route struct {
	Method     string
	Path       string
	Handlers   []Handler
	ParamNames []string
}

// Get adds the route to the router with the GET method
func (r *Router) Get(path string, handlers ...Handler) {
	r.Add(http.MethodGet, path, handlers...)
}

// Post adds the route to the router with the POST method
func (r *Router) Post(path string, handlers ...Handler) {
	r.Add(http.MethodPost, path, handlers...)
}

// Put adds the route to the router with the PUT method
func (r *Router) Put(path string, handlers ...Handler) {
	r.Add(http.MethodPut, path, handlers...)
}

// Delete adds the route to the router with the DELETE method
func (r *Router) Delete(path string, handlers ...Handler) {
	r.Add(http.MethodDelete, path, handlers...)
}

// Patch adds the route to the router with the PATCH method
func (r *Router) Patch(path string, handlers ...Handler) {
	r.Add(http.MethodPatch, path, handlers...)
}

// Head adds the route to the router with the HEAD method
func (r *Router) Head(path string, handlers ...Handler) {
	r.Add(http.MethodHead, path, handlers...)
}

// Options adds the route to the router with the OPTIONS method
func (r *Router) Options(path string, handlers ...Handler) {
	r.Add(http.MethodOptions, path, handlers...)
}

// Connect adds the route to the router with the CONNECT method
func (r *Router) Connect(path string, handlers ...Handler) {
	r.Add(http.MethodConnect, path, handlers...)
}

// Trace adds the route to the router with the TRACE method
func (r *Router) Trace(path string, handlers ...Handler) {
	r.Add(http.MethodTrace, path, handlers...)
}
