package pulse

import (
	"testing"
	"time"
)

func TestRouter_Group(t *testing.T) {
	router := NewRouter()
	api := &Group{
		Prefix: "/api",
		Router: router,
	}
	v1 := api.Group("/v1")
	v1.GET("/users", func(ctx *Context) error {
		ctx.String("users")
		return nil
	})
	v1.POST("/users/1", func(ctx *Context) error {
		ctx.String("users")
		return nil
	})
	v1.GET("/users/1", func(ctx *Context) error {
		ctx.String("users")
		return nil
	})
	v1.PUT("/users/1", func(ctx *Context) error {
		ctx.String("users")
		return nil
	})
	v1.DELETE("/users/1", func(ctx *Context) error {
		ctx.String("users")
		return nil
	})
	v1.PATCH("/users/1", func(ctx *Context) error {
		ctx.String("users")
		return nil
	})
	v1.OPTIONS("/users/1", func(ctx *Context) error {
		ctx.String("users")
		return nil
	})
	v1.HEAD("/users/1", func(ctx *Context) error {
		ctx.String("users")
		return nil
	})
	v1.Static("/users", "./static", &Static{
		Compress:      true,
		ByteRange:     false,
		IndexName:     "index.html",
		CacheDuration: 24 * time.Hour,
	})

	app.Router = router
}
