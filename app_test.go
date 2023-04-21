package pulse

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var app *Pulse

func init() {
	app = New(Config{
		AppName: "Test App",
	})
}

func TestNew(t *testing.T) {
	app := New()
	if app.config.AppName != DefaultAppName {
		t.Errorf("AppName: expected %q, actual %q", DefaultAppName, app.config.AppName)
	}
	if app.config.Network != DefaultNetwork {
		t.Errorf("Network: expected %q, actual %q", DefaultNetwork, app.config.Network)
	}

	// Test New() function with custom config
	app = New(Config{
		AppName: "Test App",
		Network: "udp",
	})
	if app.config.AppName != "Test App" {
		t.Errorf("AppName: expected %q, actual %q", "Test App", app.config.AppName)
	}
	if app.config.Network != "udp" {
		t.Errorf("Network: expected %q, actual %q", "udp", app.config.Network)
	}
}

func TestPulse_startupMessage(t *testing.T) {
	app := New(Config{
		AppName: "Test App",
	})

	addr := "localhost:8080"
	expected := "=> Server started on <" + addr + ">\n" +
		"=> App Name: " + app.config.AppName + "\n" +
		"=> Press CTRL+C to stop\n"
	actual := app.startupMessage(addr)

	if actual != expected {
		t.Errorf("startupMessage: expected %q, actual %q", expected, actual)
	}
}

func TestRouterHandler2(t *testing.T) {
	router := NewRouter()
	router.Get("/", func(ctx *Context) error {
		ctx.String("Hello, World!")
		return nil
	})

	handler := RouterHandler(router)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
}

func TestPulse_Run(t *testing.T) {
	app := New(Config{
		AppName: "test-app",
	})

	go app.Run(":9090")

	// Wait for server to start
	time.Sleep(time.Second)

	resp, err := http.Get("http://localhost:9090/")
	if err != nil {
		t.Errorf("failed to make GET request: %v", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected status code %d, got %d", http.StatusNotFound, resp.StatusCode)
	}

	err = app.Stop()
	if err != nil {
		t.Errorf("failed to stop server: %v", err)
	}
}

func TestPulse_Stop(t *testing.T) {
	app := New()

	go app.Run(":9090")

	// Wait for server to start
	time.Sleep(time.Second)

	err := app.Stop()
	if err != nil {
		t.Errorf("failed to stop server: %v", err)
	}

	// Make sure server is stopped by attempting to make a GET request
	_, err = http.Get("http://localhost:9090/")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}
