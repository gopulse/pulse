package routing

import (
	"github.com/valyala/fasthttp"
	"testing"
)

func TestRouterHandler(t *testing.T) {
	router := New()

	router.Get("/users/:id/:name", func(ctx *Context) error {
		param := ctx.Param("name")
		ctx.String(param)
		ctx.String("hello")
		return nil
	})

	err := fasthttp.ListenAndServe(":8083", RouterHandler(router))
	if err != nil {
		return
	}
}
