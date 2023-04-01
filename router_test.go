package routing

import (
	"github.com/valyala/fasthttp"
	"testing"
)

func TestRouterHandler(t *testing.T) {
	router := New()

	router.Get("/users", func(ctx *Context) error {
		err := ctx.String("hello")
		if err != nil {
			return err
		}
		return nil
	})

	err := fasthttp.ListenAndServe(":8083", RouterHandler(router))
	if err != nil {
		return
	}
}
