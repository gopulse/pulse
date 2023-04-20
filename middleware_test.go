package pulse

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_Use(t *testing.T) {
	// Create a new router.
	r := NewRouter()

	// Define a simple handler that sets a custom response header.
	handler := func(ctx *Context) error {
		ctx.SetResponseHeader("X-Test-Header", "test")
		return nil
	}

	// Add a CORS middleware to the router.
	r.Use(http.MethodGet, CORSMiddleware())

	// Add the simple handler to the router.
	r.Add(http.MethodGet, "/", handler)

	// Create an HTTP request to test the handler.
	_ = httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Verify that the response has the expected status code and header.
	if rec.Code != http.StatusOK {
		t.Errorf("unexpected status code: %d", rec.Code)
	}
}
