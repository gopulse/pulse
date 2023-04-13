package pulse

import "fmt"

type Group struct {
	Prefix string
	Router *Router
}

func (g *Group) Group(prefix string) *Group {
	return &Group{
		Prefix: g.Prefix + prefix,
		Router: g.Router,
	}
}

func (g *Group) Use(middleware Middleware) {
	g.Router.Use(g.Prefix, middleware)
}

func (g *Group) GET(path string, handlers ...Handler) {
	fmt.Println(g.Prefix + path)
	g.Router.Get(g.Prefix+path, handlers...)
}

func (g *Group) POST(path string, handlers ...Handler) {
	g.Router.Post(g.Prefix+path, handlers...)
}

func (g *Group) PUT(path string, handlers ...Handler) {
	g.Router.Put(g.Prefix+path, handlers...)
}

func (g *Group) DELETE(path string, handlers ...Handler) {
	g.Router.Delete(g.Prefix+path, handlers...)
}

func (g *Group) PATCH(path string, handlers ...Handler) {
	g.Router.Patch(g.Prefix+path, handlers...)
}

func (g *Group) OPTIONS(path string, handlers ...Handler) {
	g.Router.Options(g.Prefix+path, handlers...)
}

func (g *Group) HEAD(path string, handlers ...Handler) {
	g.Router.Head(g.Prefix+path, handlers...)
}

func (g *Group) Static(path, root string, config *Static) {
	g.Router.Static(g.Prefix+path, root, config)
}
