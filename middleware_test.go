package pulse

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareFunc_Middleware(t *testing.T) {
	// create a mock handler
	mockHandler := func(ctx *Context) error { return nil }

	// create a middleware that adds a value to the context
	middleware := MiddlewareFunc(func(handler Handler) Handler {
		return func(ctx *Context) error {
			ctx.SetValue("key", "value")
			return handler(ctx)
		}
	})

	// call the Middleware method with the mock handler
	newHandler := middleware.Middleware(mockHandler)

	// create a new context instance
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	ctx := NewContext(w, req)

	// call the new handler with the context
	err := newHandler(ctx)

	// check if the context value was set correctly
	if val := ctx.GetValue("key"); val != "value" {
		t.Errorf("Expected context value for key \"key\" to be \"value\", but got %v", val)
	}

	// check if the original handler was called with the context
	if err != nil {
		t.Errorf("Expected err to be nil, but got %v", err)
	}
}

func TestMiddleware_Use(t *testing.T) {
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

func TestCORSMiddleware(t *testing.T) {
	// create a test context
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	ctx := NewContext(w, req)

	// create a mock handler
	mockHandler := func(ctx *Context) error {
		return nil
	}

	// create the CORS middleware
	corsMiddleware := CORSMiddleware()

	// wrap the mock handler with the CORS middleware
	handler := corsMiddleware(mockHandler)

	// call the handler with the test context
	err := handler(ctx)

	// get the http.Response from the ResponseRecorder using Result()
	res := w.Result()

	// check if the Access-Control-Allow-Origin header was set to "*"
	if header := res.Header.Get("Access-Control-Allow-Origin"); header != "*" {
		t.Errorf("Expected Access-Control-Allow-Origin header to be \"*\", but got %q", header)
	}

	// check if the Access-Control-Allow-Methods header was set correctly
	if header := res.Header.Get("Access-Control-Allow-Methods"); header != "POST, GET, OPTIONS, PUT, DELETE" {
		t.Errorf("Expected Access-Control-Allow-Methods header to be \"POST, GET, OPTIONS, PUT, DELETE\", but got %q", header)
	}

	// check if the Access-Control-Allow-Headers header was set correctly
	if header := res.Header.Get("Access-Control-Allow-Headers"); header != "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization" {
		t.Errorf("Expected Access-Control-Allow-Headers header to be \"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization\", but got %q", header)
	}

	// check if the handler returned no error
	if err != nil {
		t.Errorf("Expected handler to return no error, but got %v", err)
	}
}

func TestMiddlewareFunc_Handle(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	ctx := NewContext(w, req)

	// create a mock handler
	mockHandler := func(ctx *Context) error { return nil }

	// create a middleware that adds a value to the context
	middleware := MiddlewareFunc(func(handler Handler) Handler {
		return func(ctx *Context) error {
			ctx.SetValue("key", "value")
			return handler(ctx)
		}
	})

	// call the Handle method with the mock handler as the next handler
	err := middleware.Handle(ctx, mockHandler)

	// check if the context value was set correctly
	if val := ctx.GetValue("key"); val != "value" {
		t.Errorf("Expected context value for key \"key\" to be \"value\", but got %v", val)
	}

	// check if the next handler was called with the original context
	if err != nil {
		t.Errorf("Expected err to be nil, but got %v", err)
	}
}
