package slovo

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kberov/gledki"
	"github.com/labstack/echo/v4"
)

type Stash = gledki.Stash

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
	g.Stash = Stash{
		"generator": "Slovo2",
		"version":   VERSION,
		"codename":  CODENAME,
	}

	return c.Render(200, "hello",
		Stash{
			"title":    "Здравейте!",
			"greeting": "Добре дошли!",
		},
	)
}

func handleNotFound(c *Context, err error) error {
	// TODO: I18N & L10N
	stash := Stash{"lang": c.StraniciArgs.Lang, "title": "Няма такава страница!"}
	if strings.Contains(err.Error(), `no rows`) {
		stash["mainMenu"] = mainMenu(c, c.StraniciArgs, stash)
		return c.Render(http.StatusNotFound, `not_found`, stash)
	}
	return err
}

func tryHandleCachedPage(c *Context) error {
	err := errors.New(`page is not cached.`)
	if !canCachePage(c) {
		return err
	}
	fullPath := filepath.Join(BinDir(), `domove`, c.StraniciArgs.Domain, `public`, cached, c.CanonicalPath())
	if FileIsReadable(fullPath) {
		data, _ := os.ReadFile(fullPath)
		c.Logger().Debugf("tryHandleCachedPage: %s", c.CanonicalPath())
		return c.HTMLBlob(http.StatusOK, data)
	}
	c.Logger().Debugf("tryHandleCachedPage: %s: %s", fullPath, err.Error())
	return err
}
