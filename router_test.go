package pulse

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func init() {
	app = New(Config{
		AppName: "Test App",
		Network: "tcp",
	})
}

func TestRouterHandler(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/users", func(ctx *Context) error {
		ctx.String("hello")
		return nil
	})
}

func TestRouter_find(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Get("/users/*", func(ctx *Context) error {
		ctx.String("hello")
		return nil
	})

	router.Find("GET", "/users/1")
}

func TestRouter_Static(t *testing.T) {
	router := NewRouter()

	app.Router = router

	router.Static("/users", "./static", &Static{
		Compress:      true,
		ByteRange:     false,
		IndexName:     "index.html",
		CacheDuration: 24 * time.Hour,
	})
}

func TestStatic_notFoundHandler(t *testing.T) {
	rec := httptest.NewRecorder()

	// Call the notFoundHandler method.
	options := &Static{}
	options.notFoundHandler(rec)

	// Verify that the response has the expected status code and headers.
	if rec.Code != http.StatusNotFound {
		t.Errorf("unexpected status code: got %d, want %d", rec.Code, http.StatusNotFound)
	}
	expectedContentType := "text/plain; charset=utf-8"
	actualContentType := rec.Header().Get("Content-Type")
	if actualContentType != expectedContentType {
		t.Errorf("unexpected content type: got %q, want %q", actualContentType, expectedContentType)
	}

	// Verify that the response body contains the expected text.
	expectedBody := "404 Not Found"
	actualBody := rec.Body.String()
	if actualBody != expectedBody {
		t.Errorf("unexpected body: got %q, want %q", actualBody, expectedBody)
	}
}

func TestStatic_PathRewrite(t *testing.T) {
	// Create a new request with the specified path.
	req := &http.Request{
		URL: &url.URL{
			Path: "/path/to/file/",
		},
	}

	// Call the PathRewrite method.
	options := &Static{IndexName: "index.html"}
	rewritten := options.PathRewrite(req)

	// Verify that the path was rewritten correctly.
	expectedRewritten := "/path/to/index.html"
	actualRewritten := string(rewritten)
	if actualRewritten != expectedRewritten {
		t.Errorf("unexpected path rewrite: got %q, want %q", actualRewritten, expectedRewritten)
	}
}
