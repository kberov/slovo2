/*
Package slovo contains code for the business logic of the application.
*/
package slovo

import (
	"net/http/cgi"
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const VERSION = "2024.01.05"
const CODENAME = "U+2C16 GLAGOLITIC CAPITAL LETTER UKU (Ⱆ)"

func initEcho(logger *log.Logger) *echo.Echo {
	e := echo.New()
	e.Logger = logger
	//e.GET("/", hello)...
	loadRoutes(e)
	return e
}

func loadRoutes(e *echo.Echo) {
	for _, route := range DefaultConfig.Routes {
		method := reflect.ValueOf(e).MethodByName(route.Method)
		params := []reflect.Value{
			reflect.ValueOf(route.Path),
			reflect.ValueOf(handlers[route.Handler]),
		}
		// e.GET("/", hello)
		method.Call(params)
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
