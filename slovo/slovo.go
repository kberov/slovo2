/*
Package slovo contains code for preparing and serving web pages for the site --
the front-end.
*/
package slovo

import (
	"net/http/cgi"

	"github.com/kberov/slovo2/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const VERSION = "2024.03.12-alpha-014"
const CODENAME = "U+2C16 GLAGOLITIC CAPITAL LETTER UKU (Ⱆ)"

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
	e.Use(PreferDomainStaticFiles)
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
	if err := cgi.Serve(initEcho(logger)); err != nil {
		logger.Fatal(err)
	}
}

// Start starts Echo in server mode.
func Start(logger *log.Logger) {
	logger.Fatal(initEcho(logger).Start(Cfg.Serve.Location))
}
