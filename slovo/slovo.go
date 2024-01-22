/*
Package slovo contains code for the business logic of the application.
*/
package slovo

import (
	"net/http/cgi"

	"github.com/kberov/slovo2/model"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

const VERSION = "2024.01.10"
const CODENAME = "U+2C16 GLAGOLITIC CAPITAL LETTER UKU (Ⱆ)"

func initEcho(logger *log.Logger) *echo.Echo {
	e := echo.New()
	e.Debug = Cfg.Debug
	e.Logger = logger
	model.Logger = logger
	model.DSN = Cfg.DB.DSN
	CfgR := Cfg.Renderer
	e.Renderer = GledkiMust(
		CfgR.TemplatesRoot,
		CfgR.Ext,
		CfgR.Tags,
		CfgR.LoadFiles,
		logger,
	)
	// Add middleware to the Echo instance
	e.Pre(middleware.RewriteWithConfig(Cfg.RewriteConfig))
	// Request ID middleware generates a unique id for a request.
	e.Use(middleware.RequestID())
	e.Static(Cfg.EchoStatic.Prefix, Cfg.EchoStatic.Root)
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
		if mfuncs := route.MiddlewareFuncs; mfuncs != nil && mfuncs[0] != "" {
			for _, funcName := range mfuncs {
				if f, ok := middlewareFuncs[funcName]; ok {
					definedMFuncs = append(definedMFuncs, f)
				}
			}
		}

		if route.Method == ANY {
			if definedMFuncs != nil {
				e.Any(route.Path, handlerFuncs[route.Handler], definedMFuncs...)
			} else {
				e.Any(route.Path, handlerFuncs[route.Handler])
			}
			continue
		}
		if definedMFuncs != nil {
			e.Add(route.Method, route.Path, handlerFuncs[route.Handler], definedMFuncs...).Name = route.Name
			continue
		}
		// otherwise simply add the route without []echo.MiddlewareFunc
		e.Add(route.Method, route.Path, handlerFuncs[route.Handler]).Name = route.Name
	}
}

func ServeCGI(logger *log.Logger) {
	logger.Debug("in slovo.ServeCGI()")
	e := initEcho(logger)
	if err := cgi.Serve(e); err != nil {
		e.Logger.Fatal(err)
	}
}

func Serve(logger *log.Logger) {
	logger.Debug("in slovo.Serve()")
	e := initEcho(logger)
	logger.Fatal(e.Start(Cfg.Serve.Location))
}
