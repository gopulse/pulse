package pulse

import (
	"errors"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"net"
	"net/http"
	"sync"
)

type (
	Pulse struct {
		config *Config
		server *http.Server
		Router *Router
		mx     sync.Mutex
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
	p.mx.Lock()
	defer p.mx.Unlock()
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
	if p.server == nil {
		return errors.New("server not running")
	}
	err := p.server.Shutdown(nil)
	if err != nil {
		return fmt.Errorf("failed to stop server: %v", err)
	}
	p.server = nil
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
