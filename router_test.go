package pulse

import (
	"testing"
	"time"
)

func init() {
	app = New(Config{
		AppName: "Test App",
		Network: "tcp",
	})
}

func TestRouterHandler(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/users", func(ctx *Context) error {
		ctx.String("hello")
		return nil
	})
}

func TestRouter_find(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/users/*", func(ctx *Context) error {
		ctx.String("hello")
		return nil
	})

	router.Find("GET", "/users/1")
}

func TestRouter_Static(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Static("/users", "./static", &Static{
		Compress:      true,
		ByteRange:     false,
		IndexName:     "index.html",
		CacheDuration: 24 * time.Hour,
	})
}
