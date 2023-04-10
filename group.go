package pulse

import "fmt"

type Group struct {
	prefix string
	router *Router
}

func (g *Group) Group(prefix string) *Group {
	return &Group{
		prefix: g.prefix + prefix,
		router: g.router,
	}
}

func (g *Group) Use(middleware Middleware) {
	g.router.Use(g.prefix, middleware)
}

func (g *Group) GET(path string, handlers ...Handler) {
	fmt.Println(g.prefix + path)
	g.router.Get(g.prefix+path, handlers...)
}

func (g *Group) POST(path string, handlers ...Handler) {
	g.router.Post(g.prefix+path, handlers...)
}

func (g *Group) PUT(path string, handlers ...Handler) {
	g.router.Put(g.prefix+path, handlers...)
}

func (g *Group) DELETE(path string, handlers ...Handler) {
	g.router.Delete(g.prefix+path, handlers...)
}

func (g *Group) PATCH(path string, handlers ...Handler) {
	g.router.Patch(g.prefix+path, handlers...)
}

func (g *Group) OPTIONS(path string, handlers ...Handler) {
	g.router.Options(g.prefix+path, handlers...)
}

func (g *Group) HEAD(path string, handlers ...Handler) {
	g.router.Head(g.prefix+path, handlers...)
}

func (g *Group) Static(path, root string, config *Static) {
	g.router.Static(g.prefix+path, root, config)
}
