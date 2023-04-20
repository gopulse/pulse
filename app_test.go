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

func TestPulse_Run_Stop(t *testing.T) {
	// Create a new Pulse instance.
	pulse := &Pulse{
		Router: NewRouter(),
		config: &Config{Network: "tcp"},
		server: &http.Server{},
	}

	// Start the server.
	address := "localhost:9000"
	go pulse.Run(address)

	// Wait for the server to start.
	time.Sleep(100 * time.Millisecond)

	// Make a test request to verify that the server is running.
	req, err := http.NewRequest("GET", "http://"+address, nil)
	if err != nil {
		t.Fatalf("unexpected error creating request: %v", err)
	}
	respRecorder := httptest.NewRecorder()
	pulse.server.Handler.ServeHTTP(respRecorder, req)

	err = pulse.Stop()
	if err != nil {
		t.Fatalf("failed to stop server: %v", err)
	}

	// Wait for the server to stop.
	time.Sleep(3 * time.Second)

	// Close the test request to ensure that all active connections are closed.
	err = respRecorder.Result().Body.Close()
	if err != nil {
		t.Fatalf("failed to close response body: %v", err)
	}
}
