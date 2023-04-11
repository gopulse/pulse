package pulse

import (
	"testing"
	"time"
)

var app *Pulse

func init() {
	app = New(Config{
		AppName: "Test App",
	})
}

func TestPulse_Run(t *testing.T) {
	addr := "127.0.0.1:8082"
	app.Run(addr)

	go func() {
		time.Sleep(1000 * time.Millisecond)
		app.Stop()
	}()
}
