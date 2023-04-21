package pulse

import (
	"context"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"net"
	"net/http"
	"time"
)

type (
	Pulse struct {
		config *Config
		server *http.Server
		Router *Router
	}

	Config struct {
		// AppName is the name of the app
		AppName string `json:"app_name"`

		// Network is the network to use
		Network string `json:"network"`
	}
)

const (
	// DefaultAppName is the default app name
	DefaultAppName = "Pulse"

	// DefaultNetwork is the default network
	DefaultNetwork = "tcp"
)

func New(config ...Config) *Pulse {
	app := &Pulse{
		config: &Config{},
		server: &http.Server{},
		Router: NewRouter(),
	}

	if len(config) > 0 {
		app.config = &config[0]
	}

	if app.config.AppName == "" {
		app.config.AppName = DefaultAppName
	}

	if app.config.Network == "" {
		app.config.Network = DefaultNetwork
	}

	return app
}

func (p *Pulse) Run(address string) {
	// setup handler
	handler := RouterHandler(p.Router)
	p.server.Handler = handler

	// setup listener
	listener, err := net.Listen(p.config.Network, address)
	if err != nil {
		panic(fmt.Errorf("failed to listen: %v", err))
	}

	// print startup message
	fmt.Println(p.startupMessage(listener.Addr().String()))

	// start server
	err = p.server.Serve(listener)
	if err != nil {
		fmt.Errorf("failed to start server on %s: %v", listener.Addr().String(), err)
	}
}

func (p *Pulse) Stop() error {
	// Check if the server is already stopped.
	if p.server == nil {
		return nil
	}

	// Disable HTTP keep-alive connections to prevent the server from
	// accepting any new requests.
	p.server.SetKeepAlivesEnabled(false)

	// Shutdown the server gracefully to allow existing connections to finish.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := p.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("failed to shut down server: %v", err)
	}

	// Set the server to a new instance of http.Server to allow starting it again.
	p.server = &http.Server{}

	return nil
}

func (p *Pulse) startupMessage(addr string) string {
	myFigure := figure.NewFigure("PULSE", "", true)
	myFigure.Print()

	var textOne = "=> Server started on <%s>" + "\n"
	var textTwo = "=> App Name: %s" + "\n"
	var textThree = "=> Press CTRL+C to stop" + "\n"

	return fmt.Sprintf(textOne, addr) + fmt.Sprintf(textTwo, p.config.AppName) + fmt.Sprintf(textThree)
}
