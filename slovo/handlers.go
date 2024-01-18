package slovo

import (
	"fmt"
	"net/http"

	"github.com/kberov/gledki"
	"github.com/labstack/echo/v4"
)

type Map = gledki.Stash

var spf = fmt.Sprintf

// This file contains the controllers (http handler functions) for slovo

// GET /v2/ebookform text/html
// Display

/*
TODO:
Personalize an epub file for personal usage.
Send back a link to the file to be downloaded and a password for opening
the file.
POST /v2/pepub
c.FormValue("name") - string  "First Last"
c.FormValue("email") - string "em@site.com"
c.FormValue("payed") - bool "yes|1"/"no|0"
*/
func pepubcpu(c echo.Context) error {
	return c.String(http.StatusOK, "TODO!")
}

// GET / hello
func hello(c echo.Context) error {
	c.Logger().Debugf("in hello")
	// We can use all methods of gledki.Gledki
	g := c.Echo().Renderer.(*EchoRenderer)
	g.Stash = Map{
		"generator": "Slovo2",
		"version":   VERSION,
		"codename":  CODENAME,
	}

	return c.Render(200, "hello",
		Map{
			"title":    "Здравейте!",
			"greeting": "Добре дошли!",
		},
	)
}
