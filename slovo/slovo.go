/*
Package slovo contains code for the business logic of the application.
*/
package slovo

import (
	"net/http/cgi"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const VERSION = "2024.01.10"
const CODENAME = "U+2C16 GLAGOLITIC CAPITAL LETTER UKU (Ⱆ)"

func initEcho(logger *log.Logger) *echo.Echo {
	e := echo.New()
	e.Logger = logger
	e.Renderer = GledkiMust(
		DefaultConfig.Renderer.TemplatesRoot,
		DefaultConfig.Renderer.Ext,
		DefaultConfig.Renderer.Tags,
		DefaultConfig.Renderer.LoadFiles,
		logger,
	)
	e.Static(DefaultConfig.EchoStatic.Prefix, DefaultConfig.EchoStatic.Root)

	//e.GET("/", hello)...
	loadRoutes(e)
	return e
}

// Add routes, specified in DefaultConfig.Routes to echo routes handler.
func loadRoutes(e *echo.Echo) {
	for _, route := range DefaultConfig.Routes {
		// find middleware and attach to the route if specified in configuration
		if mfuncs := route.MiddlewareFuncs; mfuncs != nil && mfuncs[0] != "" {
			var definedMFuncs []echo.MiddlewareFunc
			for _, funcName := range mfuncs {
				if f, ok := middlewareFuncs[funcName]; ok {
					definedMFuncs = append(definedMFuncs, f)
				}
			}
			e.Add(route.Method, route.Path, handlerFuncs[route.Handler], definedMFuncs...)
			continue
		}
		// otherwise simply add the route without []echo.MiddlewareFunc
		e.Add(route.Method, route.Path, handlerFuncs[route.Handler])
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
	logger.Fatal(e.Start(DefaultConfig.Serve.Location))
}
