package slovo

import (
	"net/http"
	"net/http/cgi"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const VERSION = "2024.01.05"
const CODENAME = "U+2C16 GLAGOLITIC CAPITAL LETTER UKU (Ⱆ)"

type Config struct {
	Debug      bool
	ConfigFile string
	Serve      ServeConfig
	ServeCGI   ServeCGIConfig
}

type ServeConfig struct {
	Location string
}

type ServeCGIConfig struct {
	// Default REQUEST_METHOD for testing CGI mode.
	REQUEST_METHOD  string
	SERVER_PROTOCOL string
	REQUEST_URI     string
}

var DefaultConfig Config

func init() {
	DefaultConfig = Config{
		Debug:      false,
		ConfigFile: "etc/config.yaml",
		Serve:      ServeConfig{Location: "localhost:3000"},
		ServeCGI: ServeCGIConfig{
			REQUEST_METHOD:  http.MethodGet,
			SERVER_PROTOCOL: "HTTP/1.1",
			REQUEST_URI:     "/",
		},
	}
}

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
