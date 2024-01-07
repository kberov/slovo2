package slovo

import (
	"errors"
	"net/http"

	"github.com/kberov/gledki"
	"github.com/labstack/echo/v4"
)

// This file contains the controllers (http handler functions) for slovo

// GET /v2/ebookform text/html
// Display

/*
Personalize a pdf file using pdfcpu for personal usage.
Send back a link to the file to be downnloaded and a password for opening
the file.
Store a record with order number and the link for download
snd back a message with the order num to the user saying that with this
order the file can be rtreived again in case of loss or file corruption.

POST /v2/ppdfcpu application/json
c.FormValue("name") - string  "First Last"
c.FormValue("email") - string "em@site.com"
c.FormValue("order") - string "1JK123"
*/
func ppdfcpu(c echo.Context) error {
	pdfMsg := pdfcpuMessage{Name: "Jon Смит", Email: "ala@bala.bg",
		Msg: `Gotowo: <a href://"site.com/book.pdf">Ime na kniga</a>`}
	return c.JSON(http.StatusCreated, pdfMsg)
}

// pdfcpuMessage is produced by ppdfcpu()
type pdfcpuMessage struct {
	Name  string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
	Msg   string `json:"msg" xml:"msg" form:"msg" query:"msg"`
}

/*
Displays a HTML form for ppdfcpu() and caches it on disk for subsequent static
rendering by slovo2 or Apache (when running in CGI mode).
GET /v2/epub
*/
func ppdfcpuForm(c echo.Context) error {
	return errors.New("temporary until we have a renderer")
}

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
	// We can use all methods of gledki.Gledki
	g := c.Echo().Renderer.(*EchoRenderer)
	g.Stash = gledki.Stash{
		"generator": "Slovo2",
		"version":   VERSION,
		"codename":  CODENAME,
	}

	return c.Render(200, "hello",
		gledki.Stash{
			"title":    "Здравейте!",
			"greeting": "Добре дошли!",
		},
	)
	//return c.String(http.StatusOK, "Hello, World!")
}

// We need this map because the function names are stored in yaml config as
// strings. This map is used in loadRoutes() to match HTTP handlers by name.
var handlers = map[string]echo.HandlerFunc{
	"hello":       hello,
	"ppdfcpu":     ppdfcpu,
	"ppdfcpuForm": ppdfcpuForm,
}
