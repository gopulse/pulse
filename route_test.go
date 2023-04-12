package pulse

import "testing"

func TestRoute_Get(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/users/*", func(ctx *Context) error {
		ctx.String("hello")
		return nil
	})
}

func TestRoute_Post(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Post("/users/*", func(ctx *Context) error {
		ctx.String("hello")
		return nil
	})
}

func TestRoute_Put(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Put("/users/*", func(ctx *Context) error {
		ctx.String("hello")
		return nil
	})
}

func TestRoute_Delete(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Delete("/users/*", func(ctx *Context) error {
		ctx.String("hello")
		return nil
	})
}

func TestRoute_Patch(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Patch("/users/*", func(ctx *Context) error {
		ctx.String("hello")
		return nil
	})
}

func TestRoute_Head(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Head("/users/*", func(ctx *Context) error {
		ctx.String("hello")
		return nil
	})
}

func TestRoute_Options(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Options("/users/*", func(ctx *Context) error {
		ctx.String("hello")
		return nil
	})
}

func TestRoute_Connect(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Connect("/users/*", func(ctx *Context) error {
		ctx.String("hello")
		return nil
	})
}

func TestRoute_Trace(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Trace("/users/*", func(ctx *Context) error {
		ctx.String("hello")
		return nil
	})
}
