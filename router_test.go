package routing

import (
	"github.com/valyala/fasthttp"
	"testing"
)

func TestRouterHandler(t *testing.T) {
	router := New()

	router.Post("/users/:id/:name", func(ctx *Context) error {
		param := ctx.Param("name")
		ctx.String(param)
		ctx.String("hello")
		return nil
	})

	router.Use("POST", CORSMiddleware())

	err := fasthttp.ListenAndServe(":8083", RouterHandler(router))
	if err != nil {
		return
	}
}
