package routing

import (
	"testing"
)

func TestParse(t *testing.T) {
	router := New()

	router.Get("/users/:id", func(ctx *Context) error {
		err := ctx.String("hello")
		if err != nil {
			return err
		}
		return nil
	})
}
