package pulse

import (
	"testing"
)

func TestMiddleware_Use(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/", func(ctx *Context) error {
		return nil
	})

	router.Use("GET", CORSMiddleware())
}
