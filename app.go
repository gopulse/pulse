package pulse

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/valyala/fasthttp"
	"net"
)

type (
	Pulse struct {
		config *Config
		server *fasthttp.Server
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
		server: &fasthttp.Server{},
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

func (f *Pulse) Run(address string) {
	handler := RouterHandler(f.Router)
	f.server.Handler = handler

	// setup listener
	listener, err := net.Listen(f.config.Network, address)
	if err != nil {
		panic(fmt.Errorf("failed to listen: %v", err))
	}

	// print startup message
	fmt.Println(f.startupMessage(listener.Addr().String()))

	// start server
	err = f.server.Serve(listener)
	if err != nil {
		panic(fmt.Errorf("failed to serve: %v", err))
	}
}

func (f *Pulse) Stop() {
	err := f.server.Shutdown()
	if err != nil {
		return
	}
}

func (f *Pulse) startupMessage(addr string) string {
	myFigure := figure.NewFigure("PULSE", "", true)
	myFigure.Print()

	var textOne = "=> Server started on <%s>" + "\n"
	var textTwo = "=> App Name: %s" + "\n"
	var textThree = "=> Press CTRL+C to stop" + "\n"

	return fmt.Sprintf(textOne, addr) + fmt.Sprintf(textTwo, f.config.AppName) + fmt.Sprintf(textThree)
}
