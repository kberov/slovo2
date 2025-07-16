package slovo

import (
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ANY is an aggregate for any http method.
const ANY = "ANY"

// SLOG is a regular expression capturing group to match what is possible to
// have between two slashes in an URL path. Used in RegexRules for rewriting
// urls for the Routes parser. At least three any unicode letter, dash or
// underscore.
// Note! REQUEST_URI is url-escaped at this time. We currently use Skipper to
// unnescape the raw RequestURI.
const SLOG = `([\pL\-_\d]{3,})`

// LNG is a regular expression for language notation.
const LNG = `((?:[a-z]{2}-[a-z]{2})|[a-z]{2})`

const format = `html`

// EXT is a regular expression for the requested default format.
const EXT = `(html?)`

// QS stands for QUERY_STRING - this is the rest of the URL. We match anything.
const QS = `(.*)?`

const rootAlias = `коренъ`
const guestID = 2

// Config is the root structure of the configuration for slovo. We preserve the
// case and style of each node and scalar item between YAML file and Go source
// code for best recognition.
//
//lint:file-ignore ST1003 ALL_CAPS match the ENV variable names
type Config struct {
	// Langs is a list of supported languages. The last is the default.
	Langs []string `yaml:"Langs"`
	Debug bool     `yaml:"Debug"`
	// DomoveRoot is the root folder for static content and templates of parked
	// domains. Defaults to `BinDir()/domove`.
	DomoveRoot string `yaml:"DomoveRoot"`
	// DomovePrefixes is a common list of prefixes, used in parked domains to
	// create subdomains for various purposes. These prefixes are cut and then
	// the domain name is used for domain root.
	DomovePrefixes []string `yaml:DomovePrefixes`
	// GuestID is the id of the user when a guest (unauthenticated) visits the
	// site.
	GuestID int32 `yaml:"GuestID"`
	// File is the path to the configuration file in which this structure will
	// be dumped in YAML format.
	File     string   `yaml:"File"`
	Serve    Serve    `yaml:"Serve"`
	StartCGI ServeCGI `yaml:"ServeCGI"`
	// List of routes to be created by Echo
	Routes Routes `yaml:"Routes"`
	// Arguments for instantiating GledkiRenderer
	Renderer Renderer `yaml:"Renderer"`
	// StaticRoutes - folders` routes to be served by `echo` from
	// within the installation public folder. Common for all domains. If a
	// file is not found in the domain public folder, it will fall back to
	// be served from here if found. For example request to /css/site.css
	// will be served from public/css/site.css.
	// `e.Static("/css","public/css").`
	StaticRoutes []StaticRoute `yaml:"StaticRoutes"`
	// DomoveStaticFiles is a regex of file extensions. If a file, matching the
	// regex is requested, we will look into the domain specific public folder
	// and serve it if found.
	DomoveStaticFiles string `yaml:DomoveStaticFiles`
	// DB is for database configuration. For now we use sqlite3.
	DB DBConfig `yaml:"DB"`
	// Rewrite is used to pass configuration values to
	// [middleware.RewriteWithConfig].
	Rewrite Rewrite `yaml:"Rewrite"`
	// CachePages - should we cache pages or not. true means `yes` and false
	// means `no`
	CachePages bool `yaml:"CachePages"`
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

// Renderer contains arguments, passed to GledkiMust.
type Renderer struct {
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

// Serve has configuration properties, passed to [slovo.Start].
type Serve struct {
	// Location is descrived as f.q.d.n:port
	Location string `yaml:"Location"`
}

// ServeCGI contains minimum ENV values for emulating a CGI request on
// the command line. All of these can be overridden via flags.
// See https://www.rfc-editor.org/rfc/rfc3875
type ServeCGI struct {
	HTTP_HOST      string `yaml:"HTTP_HOST"`
	REQUEST_URI    string `yaml:"REQUEST_URI"`
	REQUEST_METHOD string `yaml:"REQUEST_METHOD"`
	// SERVER_PROTOCOL used in CGI environment - HTTP/1.1. Recuired variable by
	// the cgi Go module.
	SERVER_PROTOCOL     string `yaml:"SERVER_PROTOCOL"`
	HTTP_ACCEPT_CHARSET string `yaml:"HTTP_ACCEPT_CHARSET"`
	CONTENT_TYPE        string `yaml:"CONTENT_TYPE"`
}

// Rewrite is used to pass configuration values to
// [middleware.RewriteWithConfig].
type Rewrite struct {
	SkipperFuncName string `yaml:"SkipperFuncName"`
	// Rules is a map of string to string in which the key is a regular
	// expression and the value is the resulting route mapping description
	Rules map[string]string `yaml:"Rules"`
}

/*
ToRewriteRules converts [slovo.Rewrite] to [middleware.RewriteConfig]
structure, suitable for passing to [middleware.RewriteWithConfig] and returns
it. Particularly SkipperFuncName is converted to
[middleware.RewriteConfig.Skipper] and Rules are converted to
[middleware.RewriteConfig.RegexRules]. This is somehow easier than implementing
MarshalYAML and UnmarshalYAML.
*/
func (rc Rewrite) ToRewriteRules() (rewriteConfig middleware.RewriteConfig) {
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
		/*
			  	 req.RequestURI is used by middleware#rewriteURL, but in CGI
				 environment it seems to be empty. So here we populate it
				 from URL.Path. And we do it also in server mode because
				 RequestURI is still escaped and cannot match any of our
				 regexes. We add also the RawQuery as it is needed for
				 paging of stranici and celini, and probably other items from
				 the database in the future.
				 c.Logger().Debugf("os.Getenv(REQUEST_URI): %#v", os.Getenv(`REQUEST_URI`))
				 c.Logger().Debugf("c.Request().RequestURI: %#v", c.Request().RequestURI)
				 c.Logger().Debugf(" c.Request().URL.RawQuery: %#v", c.Request().URL.RawQuery)
		*/
		var uri strings.Builder
		uri.WriteString(c.Request().URL.Path)
		if c.Request().URL.RawQuery != "" {
			uri.WriteString(`?` + c.Request().URL.RawQuery)
		}
		c.Request().RequestURI = uri.String()
		return false
	},
}

// We need this map because the function names are stored in yaml config as
// strings. This map is used in loadRoutes() to match HTTP handlerFuncs by name.
var handlerFuncs = map[string]echo.HandlerFunc{
	"hello":           hello,
	"ppdfcpu":         ppdfcpu,
	"ppdfcpuForm":     ppdfcpuForm,
	"straniciExecute": straniciExecute,
	"celiniExecute":   celiniExecute,
}

// We need this map because the function names are stored in yaml config as
// strings. These are functions only for the corresponding HandlerFunc where
// their key-names are mentioned.
var middlewareFuncs = map[string]echo.MiddlewareFunc{
	"SlovoContext": SlovoContext,
	"CachePages":   middleware.BodyDump(cachePages),
}

var defaultHost = "dev.xn--b1arjbl.xn--90ae"

// Cfg is the global configuration structure for slovo. The default is
// hardcodded and it can be dumped to YAML by using the command `slovo2 config
// dump`. To read automatically the YAML file on startup, the SLOVO_CONFIG
// environment variable must be set to the config file path.
var Cfg Config

func init() {
	// Default configuration
	Cfg.Langs = []string{"bg"}
	Cfg = Config{
		Debug:   true,
		GuestID: guestID,
		File:    "etc/config.yaml",
		Langs:   Cfg.Langs,
		Serve:   Serve{Location: spf("%s:3000", defaultHost)},
		StartCGI: ServeCGI{
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
			// Routes are not as powerful as in Mojolicious. We need the RewriteConfig.Rules below
			Route{Method: echo.GET, Path: "/", Handler: "straniciExecute", Name: "/"},
			Route{Method: ANY, Path: "/:stranica/:lang/:format", Handler: "straniciExecute",
				MiddlewareFuncs: []string{"SlovoContext", "CachePages"}, Name: "stranica"},
			Route{Method: ANY, Path: "/:stranica/:celina/:lang/:format", Handler: "celiniExecute",
				MiddlewareFuncs: []string{"SlovoContext", "CachePages"}, Name: "celina"},
			Route{Method: echo.GET, Path: "/v2/ppdfcpu", Handler: "ppdfcpuForm", Name: "ppdfcpu"},
			Route{Method: echo.POST, Path: "/v2/ppdfcpu", Handler: "ppdfcpu", Name: "ppdfcpuForm"},
		},
		Rewrite: Rewrite{
			SkipperFuncName: `URI2PathAndDontSkip`,
			Rules: map[string]string{
				// Root page in all domains has by default alias 'коренъ' and language
				// 'bg-bg'. Change the value of page_alias and the alias value of the page's
				// row in table 'stranici' for example to 'index' if you want your root page
				// to have alias 'index'. Also change the 'lang' here as desired.
				// Defaults:
				`^$`:                   spf("/%s/%s/%s", rootAlias, Cfg.Langs[0], format),
				`^/$`:                  spf("/%s/%s/%s", rootAlias, Cfg.Langs[0], format),
				spf(`^/index.%s`, EXT): spf("/%s/%s/%s", rootAlias, Cfg.Langs[0], format),
				// Станица	            /:stranica/:lang/:ext
				spf(`^/%s\.%s%s`, SLOG, EXT, QS):          "/$1/" + Cfg.Langs[0] + "/$2$3",
				spf(`^/%s\.%s\.%s%s`, SLOG, LNG, EXT, QS): "/$1/$2/$3$4",

				// Целина      /:stranica/:celina/:lang/:ext
				// for now we have content only in bulgarian
				spf(`^/%s/%s\.%s%s`, SLOG, SLOG, EXT, QS):          spf("/$1/$2/%s/$3$4", Cfg.Langs[0]),
				spf(`^/%s/%s\.%s\.%s%s`, SLOG, SLOG, LNG, EXT, QS): "/$1/$2/$3/$4$5",
			},
		},
		Renderer: Renderer{
			// Templates root folder. Must exist.
			TemplateRoots: []string{filepath.Join(HomeDir(), "templates")},
			Ext:           ".htm",
			// Delimiters for template tags
			Tags: [2]string{"${", "}"},
			// Should the template files be loaded at start?
			LoadFiles: false,
		},
		StaticRoutes: []StaticRoute{
			StaticRoute{Prefix: "/css", Root: "public/css"},
			StaticRoute{Prefix: "/fonts", Root: "public/fonts"},
			StaticRoute{Prefix: "/img", Root: "public/img"},
			StaticRoute{Prefix: "/js", Root: "public/js"},
		},
		DomoveStaticFiles: `(?i:\.(?:|png|webp|gif|jpe?g|js|css|html|pdf|woff2?))$`,
		DomovePrefixes:    []string{`dev.`, `www.`, `qa.`, `bg.`, `en.`},
		DB: DBConfig{
			DSN: filepath.Join(HomeDir(), "data/slovo.dev.sqlite"),
		},
		CachePages: true,
	}

	Cfg.DomoveRoot = filepath.Join(HomeDir(), `domove`)
	Cfg.CachePages = false
} // end init()
