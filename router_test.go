package pulse

import (
	"fmt"
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

	router.Get("/users/*", func(ctx *Context) error {
		ctx.String("hello")
		return nil
	})

}

func TestCORSMiddleware(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/", func(ctx *Context) error {
		return nil
	})

	router.Use("GET", CORSMiddleware())

}

func TestContext_SetCookie(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/", func(ctx *Context) error {
		cookie := Cookie{
			Name:        "Test Cookie 1",
			Value:       "Test Cookie 1",
			Path:        "/",
			Domain:      "localhost",
			MaxAge:      0,
			Expires:     time.Now().Add(24 * time.Hour),
			Secure:      false,
			HTTPOnly:    false,
			SameSite:    "Lax",
			SessionOnly: false,
		}
		ctx.SetCookie(&cookie)
		return nil
	})

}

func TestContext_GetCookie(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/", func(ctx *Context) error {
		cookie := ctx.GetCookie("test")
		ctx.String(cookie)
		return nil
	})

}

func TestContext_SetHeader(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/", func(ctx *Context) error {
		ctx.SetHeader("Test Header", "test header value")
		fmt.Println(ctx.GetHeader("test"))
		return nil
	})

}

func TestRouter_Static(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Static("/", "./static", &Static{
		Compress:      true,
		ByteRange:     false,
		IndexName:     "index.html",
		CacheDuration: 24 * time.Hour,
	})

}
