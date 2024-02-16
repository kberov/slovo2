package slovo

import (
	"net/http"
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
const SLOG = `([\pL\-_\d]+)`

// LNG is a regular expression for language notation.
const LNG = `((?:[a-z]{2}-[a-z]{2})|[a-z]{2})`

// EXT is a regular expression for the requested format.
const EXT = `(html?)`

// QS stands for QUERY_STRING - this is the rest of the URL. We match anything.
const QS = `(.*)?`

const rootPageAlias = `коренъ`

// Config is the root structure of the configuration for slovo. We preserve the
// case and style of each node and scalar item for best recognition between YAML
// file and Go source code.
//
//lint:file-ignore ST1003 ALL_CAPS match the ENV variable names
type Config struct {
	// Languages is a list of supported languages. the last is the default.
	Languages  []string       `yaml:"Languages"`
	Debug      bool           `yaml:"Debug"`
	ConfigFile string         `yaml:"ConfigFile"`
	Serve      ServeConfig    `yaml:"Serve"`
	ServeCGI   ServeCGIConfig `yaml:"ServeCGI"`
	// List of routes to be created by Echo
	Routes Routes `yaml:"Routes"`
	// Arguments for GledkiRenderer
	Renderer RendererConfig `yaml:"Renderer"`
	// Directories for static content. For example request to /css/site.css
	// will be served from public/css/site.css.
	// `e.Static("/css","public/css").`
	StaticRoutes  []StaticRoute `yaml:"StaticRoutes"`
	DB            DBConfig      `yaml:"DB"`
	RewriteConfig RewriteConfig `yaml:"RewriteConfig"`
}

type DBConfig struct {
	DSN string `yaml:"DSN"`
}

type StaticRoutes []StaticRoute

// StaticRoute describes a file path which will be served by echo.
type StaticRoute struct {
	Prefix string `yaml:"Prefix"`
	Root   string `yaml:"Root"`
}

type RendererConfig struct {
	TemplateRoots []string  `yaml:"TemplateRoots"`
	Ext           string    `yaml:"Ext"`
	Tags          [2]string `yaml:"Tags"`
	LoadFiles     bool      `yaml:"LoadFiles"`
}

type Route struct {
	// Method is a method name from echo.Echo.
	Method string `yaml:"Method"`
	// Handler stores a HTTP handler function name as string.
	// It is not possible to lookup a function by its name (as a string) in Go,
	// but we need to store the function names in the configuration file to
	// easily enable/disable a route. So we use a map in slovo/handlers.go `var
	// handlers = map[string]func(c echo.Context) error`
	Handler string `yaml:"Handler"`
	// Path is the REQUEST_PATH
	Path string `yaml:"Path"`
	// MiddlewareFuncs is optional
	MiddlewareFuncs []string `yaml:"MiddlewareFunc"`
	// Name is the name of the route. Used to generate URIs. See
	// https://echo.labstack.com/docs/routing#route-naming
	Name string `yaml:"Name"`
}

type Routes []Route

type ServeConfig struct {
	// Location is descrived as f.q.d.n:port
	Location string `yaml:"Location"`
}

// ServeCGIConfig contains minimum ENV values for emulating a CGI request on
// the command line. See https://www.rfc-editor.org/rfc/rfc3875
type ServeCGIConfig struct {
	HTTP_HOST      string `yaml:"HTTP_HOST"`
	REQUEST_METHOD string `yaml:"REQUEST_METHOD"`
	// SERVER_PROTOCOL used in CGI environment - HTTP/1.1. Recuired variable by
	// the cgi Go module.
	SERVER_PROTOCOL     string `yaml:"SERVER_PROTOCOL"`
	REQUEST_URI         string `yaml:"REQUEST_URI"`
	HTTP_ACCEPT_CHARSET string `yaml:"HTTP_ACCEPT_CHARSET"`
	CONTENT_TYPE        string `yaml:"CONTENT_TYPE"`
}

type RewriteConfig struct {
	SkipperFuncName string `yaml:"SkipperFuncName"`
	// Rules is a map of string to string in which the key is a regular
	// expression and the value is the resulting route mapping description
	Rules map[string]string `yaml:"Rules"`
}

/*
ToRewriteRules converts slovo.RewriteConfig to middleware.RewriteConfig
structure, suitable for passing to middleware.RewriteWithConfig() and returns
it. Particularly SkipperFuncName is converted to
middleware.RewriteConfig.Skipper and Rules are converted to
middleware.RewriteConfig.RegexRules. This is somehow easier than implementing
MarshalYAML and UnmarshalYAML.
*/
func (rc RewriteConfig) ToRewriteRules() (rewriteConfig middleware.RewriteConfig) {
	rewriteConfig = middleware.RewriteConfig{
		Skipper:    rewriteConfigSkippers[rc.SkipperFuncName],
		RegexRules: make(map[*regexp.Regexp]string),
	}
	for k, v := range rc.Rules {
		rewriteConfig.RegexRules[regexp.MustCompile(k)] = v
	}
	return
}

var rewriteConfigSkippers = map[string]middleware.Skipper{
	`URI2PathAndDontSkip`: func(c echo.Context) bool {
		// req.RequestURI is used by middleware#rewriteURL, but in CGI
		// environment it seems to be empty. So here we populate it
		// from URL.Path. And we do it unconditionally because
		// RequestURI is still escaped and cannot match any of our
		// regexes.
		c.Request().RequestURI = c.Request().URL.Path
		return false
	},
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
		Routes: Routes{
			// Routes are not as pawerful as in Mojolicious. We need the RewriteConfig.Rules below
			Route{Method: echo.GET, Path: "/", Handler: "straniciExecute", Name: "/"},
			Route{Method: ANY, Path: "/:stranica/:lang/:format", Handler: "straniciExecute"},
			Route{Method: ANY, Path: "/:stranica/:celina/:lang/:format", Handler: "celiniExecute"},
			Route{Method: echo.GET, Path: "/v2/ppdfcpu", Handler: "ppdfcpuForm", Name: "ppdfcpu"},
			Route{Method: echo.POST, Path: "/v2/ppdfcpu", Handler: "ppdfcpu", Name: "ppdfcpuForm"},
		},
		RewriteConfig: RewriteConfig{
			SkipperFuncName: `URI2PathAndDontSkip`,
			Rules: map[string]string{
				// Root page in all domains has by default alias 'коренъ' and language
				// 'bg-bg'. Change the value of page_alias and the alias value of the page's
				// row in table 'stranici' for example to 'index' if you want your root page
				// to have alias 'index'. Also change the 'lang' here as desired.
				// Defaults:
				`^$`:                   spf("/%s/bg/html", rootPageAlias),
				`^/$`:                  spf("/%s/bg/html", rootPageAlias),
				spf(`^/index.%s`, EXT): spf("/%s/bg/html", rootPageAlias),
				// Станица	            /:stranica/:lang/:ext
				spf(`^/%s\.%s%s`, SLOG, EXT, QS):          "/$1/bg/$2$3",
				spf(`^/%s\.%s\.%s%s`, SLOG, LNG, EXT, QS): "/$1/$2/$3$4",

				// Целина      /:stranica/:celina/:lang/:ext
				// for now we have content only in bulgarian
				spf(`^/%s/%s\.%s%s`, SLOG, SLOG, EXT, QS):          "/$1/$2/bg/$3$4",
				spf(`^/%s/%s\.%s\.%s%s`, SLOG, SLOG, LNG, EXT, QS): "/$1/$2/$3/$4$5",
			},
		},
		Renderer: RendererConfig{
			// Templates root folder. Must exist
			TemplateRoots: []string{"templates"},
			Ext:           ".htm",
			// Delimiters for template tags
			Tags: [2]string{"${", "}"},
			// Should the files be loaded at start?
			LoadFiles: false,
		},
		// Static files routes to be seved by echo.
		StaticRoutes: []StaticRoute{
			StaticRoute{Prefix: "/css", Root: "public/css"},
			StaticRoute{Prefix: "/fonts", Root: "public/fonts"},
			StaticRoute{Prefix: "/img", Root: "public/img"},
		},
		DB: DBConfig{
			DSN: "data/slovo.dev.sqlite",
		},
	}
}
