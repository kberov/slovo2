package slovo

import (
	"net/http"
	"net/http/cgi"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const VERSION = "2024.01.05"
const CODENAME = "U+2C16 GLAGOLITIC CAPITAL LETTER UKU (Ⱆ)"


func initEcho(logger *log.Logger) *echo.Echo {
	e := echo.New()
	e.Logger = logger
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	return e
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
