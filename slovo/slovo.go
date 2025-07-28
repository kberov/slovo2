/*
Package slovo contains code for preparing and serving web pages for the site --
the front-end.
*/
package slovo

import (
	"net/http/cgi"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/kberov/slovo2/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	flag "github.com/spf13/pflag"
)

const VERSION = "2024.04.11-alpha-015"
const CODENAME = "U+2C16 GLAGOLITIC CAPITAL LETTER UKU (Ⱆ)"

// DefaultLogHeader = `${prefix}:${time_rfc3339}:${level}:${short_file}:${line}`
// See possible placeholders in  github.com/labstack/gommon@v0.4.2/log/log.go function log()
// See https://echo.labstack.com/docs/customization
const DefaultLogHeader = `${prefix}:${level}:${short_file}:${line}`

// Logger is an instance of github.com/labstack/gommon/log.
var Logger *log.Logger

// Bin is the file name with which the program is run - the last element
// of the full path to it. Usually this is 'slovo2'. See [filepath.Base].
var Bin string = "slovo2"

func init() {
	Bin = filepath.Base(os.Args[0])
	Logger = log.New(Bin)
	//TODO: Add configuration for log file output.
	Logger.SetOutput(os.Stderr)
	Logger.SetHeader(DefaultLogHeader)
}

func initEcho(logger *log.Logger) *echo.Echo {
	e := echo.New()
	e.Debug = Cfg.Debug
	e.Logger = logger
	model.Logger = logger
	model.DSN = Cfg.DB.DSN
	CfgR := Cfg.Renderer
	e.Renderer = GledkiMust(
		CfgR.TemplateRoots,
		CfgR.Ext,
		CfgR.Tags,
		CfgR.LoadFiles,
		logger,
	)
	// Use our binder which embeds echo.DefaultBinder
	e.Binder = &Binder{}
	// Add middleware to the Echo instance
	e.Pre(middleware.RewriteWithConfig(Cfg.Rewrite.ToRewriteRules()))
	// Request ID middleware generates a unique id for a request.
	e.Use(preferDomainStaticFiles)
	e.Use(middleware.RequestID())
	// Add directories in which the files will be served as they are.
	for _, path := range Cfg.StaticRoutes {
		e.Static(path.Prefix, path.Root)
	}
	// TODO add Validator  and other needed stugff. See
	// https://echo.labstack.com/docs/customization
	// e.GET("/", hello)...
	loadRoutes(e)
	return e
}

// Add routes, specified in DefaultConfig.Routes to echo's routes handler. See
// https://echo.labstack.com/docs/routing
func loadRoutes(e *echo.Echo) {
	for _, route := range Cfg.Routes {
		// find middleware and attach to the route if specified in configuration
		var definedMFuncs []echo.MiddlewareFunc
		for _, funcName := range route.MiddlewareFuncs {
			// e.Logger.Debugf("route:%s;MiddlewareFunc: %s", route.Path, funcName)
			if f, ok := middlewareFuncs[funcName]; ok {
				definedMFuncs = append(definedMFuncs, f)
			}
		}
		if route.Method == ANY {
			e.Any(route.Path, handlerFuncs[route.Handler], definedMFuncs...)
			continue
		}
		e.Add(route.Method, route.Path, handlerFuncs[route.Handler], definedMFuncs...).Name = route.Name
	}
}

// StartCGI starts Echo in CGI mode.
func StartCGI(logger *log.Logger) {
	CgiInitEnvVarsFromConfig()
	if err := cgi.Serve(initEcho(logger)); err != nil {
		logger.Fatal(err)
	}
}

// Start starts Echo in server mode.
func Start(logger *log.Logger) {
	logger.Fatal(initEcho(logger).Start(Cfg.Serve.Location))
}

// CgiInitFlags sets a flagset to set minimum ENV values (with defaults from
// config file or default config) for emulating a CGI request on the command
// line. It will be used in the cgi command or independently when built as a
// separate executable out of the common command set for the application.
func CgiInitFlags(flags *flag.FlagSet) {
	flags.StringVarP(
		&Cfg.StartCGI.HTTP_HOST,
		"HTTP_HOST", "H",
		Cfg.StartCGI.HTTP_HOST, "The server host to which the client request is directed.")
	flags.StringVarP(
		&Cfg.StartCGI.REQUEST_URI,
		"REQUEST_URI", "U",
		Cfg.StartCGI.REQUEST_URI, "Request URI")
	flags.StringVarP(
		&Cfg.StartCGI.REQUEST_METHOD,
		"REQUEST_METHOD", "M",
		Cfg.StartCGI.REQUEST_METHOD, "Request method")
	flags.StringVarP(
		&Cfg.StartCGI.SERVER_PROTOCOL,
		"SERVER_PROTOCOL", "P",
		Cfg.StartCGI.SERVER_PROTOCOL, "Server protocol")
	flags.StringVarP(
		&Cfg.StartCGI.HTTP_ACCEPT_CHARSET,
		"HTTP_ACCEPT_CHARSET", "C",
		Cfg.StartCGI.HTTP_ACCEPT_CHARSET, "Accept-Charset")
	flags.StringVarP(
		&Cfg.StartCGI.CONTENT_TYPE,
		"CONTENT_TYPE", "T",
		Cfg.StartCGI.CONTENT_TYPE, "Content-Type")

}

// CgiInitEnvVarsFromConfig initialises environment variables from Cfg, if not
// provided, to emulate a CGI environment for testing on the command line.
func CgiInitEnvVarsFromConfig() {
	if os.Getenv("GATEWAY_INTERFACE") != "" {
		return
	}

	// minimum ENV values for emulating a CGI request on the command line
	var env = map[string]string{
		"GATEWAY_INTERFACE": "CGI/1.1",
		"SERVER_PROTOCOL":   Cfg.StartCGI.SERVER_PROTOCOL,
		"REQUEST_METHOD":    Cfg.StartCGI.REQUEST_METHOD,
		"HTTP_HOST":         Cfg.StartCGI.HTTP_HOST,
		//"HTTP_REFERER":        "elsewhere",
		//"HTTP_USER_AGENT":     "slovo2client",
		"HTTP_ACCEPT_CHARSET": Cfg.StartCGI.HTTP_ACCEPT_CHARSET,
		// "HTTP_FOO_BAR":    "baz",
		"REQUEST_URI": escapeRequestURI(Cfg.StartCGI.REQUEST_URI),
		// "CONTENT_LENGTH":  "123",
		"CONTENT_TYPE": Cfg.StartCGI.CONTENT_TYPE,
		// "REMOTE_ADDR":     "5.6.7.8",
		// "REMOTE_PORT":     "54321",
	}
	for k, v := range env {
		if os.Getenv(k) == "" {
			// Logger.Debugf("Setting %s: %s", k, v)
			os.Setenv(k, v)
		}
	}
}

// escapeRequestURI() escapes the passed on the commandline address as if a user
// agent did it.
func escapeRequestURI(uri string) string {
	uri, _ = strings.CutPrefix(uri, `/`)
	parts := strings.Split(uri, `/`)
	uri = ""
	for _, p := range parts {
		uri += `/` + url.PathEscape(p)
	}
	return uri
}
