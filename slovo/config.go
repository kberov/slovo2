package slovo

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

//lint:file-ignore ST1003 ALL_CAPS match the ENV variable names
type Config struct {
	Debug      bool
	ConfigFile string
	Serve      ServeConfig
	ServeCGI   ServeCGIConfig
	// List of routes to be created by Echo
	Routes []Route
	// Arguments for GledkiRenderer
	Renderer RendererConfig
	// Directory for static content. For example e.Static("/","public").
	EchoStatic EchoStaticConfig
	DB         DBConfig
}

type DBConfig struct {
	DSN    string
	tables []string
}

type EchoStaticConfig struct {
	Prefix string
	Root   string
}

type RendererConfig struct {
	TemplatesRoot string
	Ext           string
	Tags          [2]string
	LoadFiles     bool
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
	// Path is the REQUEST_PATH
	Path string
	// MiddlewareFuncs is optional
	MiddlewareFuncs []string
}

type ServeConfig struct {
	Location string
}

// ServeCGIConfig contains minimum ENV values for emulating a CGI request on
// the command line.
type ServeCGIConfig struct {
	REQUEST_METHOD  string
	SERVER_PROTOCOL string
	REQUEST_URI     string
}

var Cfg Config

func init() {
	// Default configuration
	Cfg = Config{
		Debug:      false,
		ConfigFile: "etc/config.yaml",
		Serve:      ServeConfig{Location: "localhost:3000"},
		ServeCGI: ServeCGIConfig{
			REQUEST_METHOD:  http.MethodGet,
			SERVER_PROTOCOL: "HTTP/1.1",
			REQUEST_URI:     "/",
		},
		// Store methods by names in YAML!
		Routes: []Route{
			Route{Method: echo.GET, Path: "/", Handler: "hello"},
			Route{Method: echo.GET, Path: "/v2/ppdfcpu", Handler: "ppdfcpuForm"},
			Route{Method: echo.POST, Path: "/v2/ppdfcpu", Handler: "ppdfcpu"},
		},
		Renderer: RendererConfig{
			// Templates root folder. Must exist
			TemplatesRoot: "templates",
			Ext:           ".htm",
			// Delimiters for template tags
			Tags: [2]string{"${", "}"},
			// Should the files be loaded at start?
			LoadFiles: false,
		},
		EchoStatic: EchoStaticConfig{
			Prefix: "/",
			Root:   "public",
		},
		DB: DBConfig{
			DSN:    "data/slovo.dev.sqlite",
			tables: []string{},
		},
	}
}
