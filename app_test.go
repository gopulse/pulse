package pulse

var app *Pulse

func init() {
	app = New(Config{
		AppName: "Test App",
	})
}
