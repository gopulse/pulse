package pulse

import (
	"testing"
	"time"
)

func TestContext_Param(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/users/:id", func(ctx *Context) error {
		ctx.String(ctx.Param("id"))
		return nil
	})
}

func TestContext_Query(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/users", func(ctx *Context) error {
		ctx.String(ctx.Query("name"))
		return nil
	})
}

func TestContext_Abort(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/", func(ctx *Context) error {
		ctx.Abort()
		return nil
	})
}

func TestContext_String(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/", func(ctx *Context) error {
		ctx.String("Test String")
		return nil
	})
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

func TestContext_ClearCookie(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/", func(ctx *Context) error {
		ctx.ClearCookie("test")
		return nil
	})
}

func TestContext_SetHeader(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/", func(ctx *Context) error {
		ctx.SetResponseHeader("Test Header", "test header value")
		return nil
	})

}

func TestContext_GetHeader(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/", func(ctx *Context) error {
		ctx.String(ctx.GetResponseHeader("test"))
		return nil
	})
}

func TestContext_SetData(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/", func(ctx *Context) error {
		ctx.SetData("test", "test data")
		return nil
	})
}

func TestContext_GetData(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/", func(ctx *Context) error {
		ctx.String(ctx.GetData("test"))
		return nil
	})
}

func TestContext_Next(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/", func(ctx *Context) error {
		err := ctx.Next()
		if err != nil {
			return err
		}
		return nil
	})
}

func TestContext_Status(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/", func(ctx *Context) error {
		ctx.Status(200)
		return nil
	})
}

func TestContext_JSON(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/", func(ctx *Context) error {
		ctx.JSON(200, map[string]string{"test": "test"})
		return nil
	})

}

func TestContext_Accepts(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/users", func(ctx *Context) error {
		ctx.SetRequestHeader("Accept", "application/json")

		accepts := ctx.Accepts("application/json", "text/html")

		if accepts == "application/json" {
			ctx.JSON(200, map[string]string{"test": "test"})
		} else {
			ctx.String("text/html")
		}
		return nil
	})

	app.Run(":8083")
}
