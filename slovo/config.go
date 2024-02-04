package slovo

import (
	"net/http"
	"net/url"
	"regexp"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ANY is an aggregata for any http method.
const ANY = "ANY"

// SLOG is a regular expression capturing group to match what is possible to
// have between two slashes in an URL path. Used in RegexRules for rewriting
// urls for the Routes parser. At least two and up to ten any unicode
// letter, dash or underscore.
// Note! REQUEST_URI is url-escaped at this time. We currently use Skipper to
// unnescape the raw RequestURI.
const SLOG = `([\pL\-_]+)`

// LNG is a regular expression for language notation.
const LNG = `((?:[a-z]{2}-[a-z]{2})|[a-z]{2})`

// EXT is a regular expression for the requested format.
const EXT = `(html?)`

// QS stands for QUERY_STRING - this is the rest of the URL. We match anything.
const QS = `(.*)?`

//lint:file-ignore ST1003 ALL_CAPS match the ENV variable names
type Config struct {
	// Languages is a list of supported languages. the last is the default.
	Languages  []string
	Debug      bool
	ConfigFile string
	Serve      ServeConfig
	ServeCGI   ServeCGIConfig
	// List of routes to be created by Echo
	Routes []Route
	// Arguments for GledkiRenderer
	Renderer RendererConfig
	// Directory for static content. For example e.Static("/","public").
	EchoStatic    EchoStaticConfig
	DB            DBConfig
	RewriteConfig middleware.RewriteConfig
}

type DBConfig struct {
	DSN    string
	Tables []string
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
	// Name is the name of the route. Used to generate URIs. See
	// https://echo.labstack.com/docs/routing#route-naming
	Name string
}

type ServeConfig struct {
	Location string
}

// ServeCGIConfig contains minimum ENV values for emulating a CGI request on
// the command line. See https://www.rfc-editor.org/rfc/rfc3875
type ServeCGIConfig struct {
	HTTP_HOST           string
	REQUEST_METHOD      string
	SERVER_PROTOCOL     string
	REQUEST_URI         string
	HTTP_ACCEPT_CHARSET string
	CONTENT_TYPE        string
}

var Cfg Config

// We need this map because the function names are stored in yaml config as
// strings. This map is used in loadRoutes() to match HTTP handlerFuncs by name.
var handlerFuncs = map[string]echo.HandlerFunc{
	"hello":           hello,
	"ppdfcpu":         ppdfcpu,
	"ppdfcpuForm":     ppdfcpuForm,
	"straniciExecute": straniciExecute,
	"celiniExecute":   celiniExecute,
}

// This map is for the same purpose as above but for one or more middleware
// functions for the corresponding HandlerFunc.
var middlewareFuncs = map[string]echo.MiddlewareFunc{}
var defaultHost = "dev.xn--b1arjbl.xn--90ae"

func init() {
	// Default configuration
	Cfg = Config{
		Languages:  []string{"bg"},
		Debug:      true,
		ConfigFile: "etc/config.yaml",
		Serve:      ServeConfig{Location: spf("%s:3000", defaultHost)},
		ServeCGI: ServeCGIConfig{
			// These are set as environment variables when the command `cgi` is
			// executed on the command line and if they are not passed as flags
			// or not set by the environment. These are the default values for
			// flags in command `cgi`.
			HTTP_HOST:           defaultHost,
			REQUEST_METHOD:      http.MethodGet,
			SERVER_PROTOCOL:     "HTTP/1.1",
			REQUEST_URI:         "/",
			HTTP_ACCEPT_CHARSET: "utf-8",
			CONTENT_TYPE:        "text/html",
		},
		// Store methods by names in YAML!
		Routes: []Route{
			// Routes are not as pawerful as in Mojolicious. We need the RegexRules below
			Route{Method: echo.GET, Path: "/", Handler: "straniciExecute", Name: "/"},
			Route{Method: ANY, Path: "/:stranica/:lang/:format", Handler: "straniciExecute"},
			// any  => '/<page_alias:str>/<paragraph_alias:cel>.<lang:lng>.html',
			Route{Method: ANY, Path: "/:stranica/:celina/:lang/:format", Handler: "celiniExecute"},
			Route{Method: echo.GET, Path: "/v2/ppdfcpu", Handler: "ppdfcpuForm", Name: "ppdfcpu"},
			Route{Method: echo.POST, Path: "/v2/ppdfcpu", Handler: "ppdfcpu", Name: "ppdfcpuForm"},
		},
		RewriteConfig: middleware.RewriteConfig{
			// TODO: think how to assign this function when parsing yaml. We
			// need some custom unmarshaller.
			Skipper: func(c echo.Context) bool {
				req := c.Request()
				escapedURI, err := url.PathUnescape(req.RequestURI)
				if err != nil {
					c.Logger().Warnf("error unescaping path: %w", err)
					return false
				}
				req.RequestURI = escapedURI
				return false
			},
			RegexRules: map[*regexp.Regexp]string{
				// Root page in all domains has by default alias 'коренъ' and language
				// 'bg-bg'. Change the value of page_alias and the alias value of the page's
				// row in table 'stranici' for example to 'index' if you want your root page
				// to have alias 'index'. Also change the 'lang' here as desired.
				// Defaults:
				regexp.MustCompile("^$"):                    "/коренъ/bg/html",
				regexp.MustCompile("^/$"):                   "/коренъ/bg/html",
				regexp.MustCompile(spf("^/index.%s$", EXT)): "/коренъ/bg/html",
				// Страница	            /:stranica/:lang/:ext
				regexp.MustCompile(spf(`^/%s\.%s%s`, SLOG, EXT, QS)):          "/$1/bg/$2$3",
				regexp.MustCompile(spf(`^/%s\.%s\.%s%s`, SLOG, LNG, EXT, QS)): "/$1/$2/$3$4",
				// Целина      /:stranica/:celina/:lang/:ext
				// for now we have content only in bulgarian
				regexp.MustCompile(spf(`^/%s/%s\.%s%s`, SLOG, SLOG, EXT, QS)):          "/$1/$2/bg/$3$4",
				regexp.MustCompile(spf(`^/%s/%s\.%s\.%s%s`, SLOG, SLOG, LNG, EXT, QS)): "/$1/$2/$3/$4$5",
			},
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
			DSN: "data/slovo.dev.sqlite",
			// This may not be needed in the Go implementation - not used for
			// now, as here the implementation is more static.
			Tables: []string{"domove", "stranici", "celini", "products"},
		},
	}
}
