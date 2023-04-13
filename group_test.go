package pulse

import "testing"

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

	app.Router = router
}
