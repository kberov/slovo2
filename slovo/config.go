package slovo

import (
	"net/http"
)

//lint:file-ignore ST1003 ALL_CAPS match the ENV variable names
type Config struct {
	Debug      bool
	ConfigFile string
	Serve      ServeConfig
	ServeCGI   ServeCGIConfig
	Routes     []Route
}

type Route struct {
	// Method is a method name from echo.Echo.
	Method string
	// Handler stores a HTTP handler name as string.
	// It is not possible to lookup a function by its name (as a string) in Go,
	// but we need to store the function names in the configuration file to
	// easily enable/disable a route. So we use a map in slovo/handlers.go `var
	// handlers = map[string]func(c echo.Context) error`
	Handler string
	Path    string
}

type ServeConfig struct {
	Location string
}

// ServeCGIConfig contains minimum ENV values for emulating a CGI request on
// the command line
type ServeCGIConfig struct {
	REQUEST_METHOD  string
	SERVER_PROTOCOL string
	REQUEST_URI     string
}

var DefaultConfig Config

func init() {
	DefaultConfig = Config{
		Debug:      false,
		ConfigFile: "etc/config.yaml",
		Serve:      ServeConfig{Location: "localhost:3000"},
		ServeCGI: ServeCGIConfig{
			REQUEST_METHOD:  http.MethodGet,
			SERVER_PROTOCOL: "HTTP/1.1",
			REQUEST_URI:     "/",
		},
		// How to store methods by names in YAML
		Routes: []Route{
			Route{Method: http.MethodGet, Path: "/", Handler: "hello"},
			Route{Method: http.MethodPost, Path: "/v2/ppdfcpu", Handler: "ppdfcpu"},
		},
	}
}
