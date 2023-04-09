package routing

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"testing"
	"time"
)

func TestRouterHandler(t *testing.T) {
	router := New()

	router.Get("/users/*", func(ctx *Context) error {
		ctx.String("hello")
		return nil
	})

	err := fasthttp.ListenAndServe(":8083", RouterHandler(router))
	if err != nil {
		return
	}
}

func TestCORSMiddleware(t *testing.T) {
	router := New()

	router.Get("/", func(ctx *Context) error {
		return nil
	})

	router.Use("GET", CORSMiddleware())

	err := fasthttp.ListenAndServe(":8083", RouterHandler(router))
	if err != nil {
		return
	}
}

func TestContext_SetCookie(t *testing.T) {
	router := New()

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

	err := fasthttp.ListenAndServe(":8083", RouterHandler(router))
	if err != nil {
		return
	}
}

func TestContext_GetCookie(t *testing.T) {
	router := New()

	router.Get("/", func(ctx *Context) error {
		cookie := ctx.GetCookie("test")
		ctx.String(cookie)
		return nil
	})

	err := fasthttp.ListenAndServe(":8083", RouterHandler(router))
	if err != nil {
		return
	}
}

func TestContext_SetHeader(t *testing.T) {
	router := New()

	router.Get("/", func(ctx *Context) error {
		ctx.SetHeader("Test Header", "test header value")
		fmt.Println(ctx.GetHeader("test"))
		return nil
	})

	err := fasthttp.ListenAndServe(":8083", RouterHandler(router))
	if err != nil {
		return
	}
}

func TestRouter_Static(t *testing.T) {
	router := New()

	router.Static("/", "./static", &Static{
		Compress:      true,
		ByteRange:     false,
		IndexName:     "index.html",
		CacheDuration: 24 * time.Hour,
	})

	err := fasthttp.ListenAndServe(":8083", RouterHandler(router))
	if err != nil {
		return
	}
}
