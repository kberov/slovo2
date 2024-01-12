package slovo

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/kberov/gledki"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
)

var spf = fmt.Sprintf

// This file contains the controllers (http handler functions) for slovo

// GET /v2/ebookform text/html
// Display

/*
Personalizes a pdf file using pdfcpu for personal usage.
Sends back a link to the file to be downnloaded and a password for opening
the file.
Store a record with order number and the link for download.
Sends back a message with the order num to the user saying that with this
order the file can be retreived again in case of loss or file corruption.

POST /v2/ppdfcpu application/json
c.FormValue("name") - string  "First Last"
c.FormValue("email") - string "em@site.com"
c.FormValue("order") - string "1JK123"
*/
func ppdfcpu(c echo.Context) error {
	passw := random.String(8, random.Uppercase, random.Numeric)
	// TODO: store it in the yaml struct for the book and get it from there.
	t := time.Now()
	year, month, day := t.Date()
	hour, min := t.Hour(), t.Minute()

	sourceFileName := "data/pdf/IS.pdf"
	dir := spf("data/pdf/%d/%d/%d", year, month, day)
	downloadFileName := spf("%s/%d%d-%sIS.pdf", dir, hour, min, passw[:3])
	pdfMsg := pdfcpuMessage{
		Download: downloadFileName,
		Passwd:   passw,
		Name:     "Jon Смит",
		Email:    "ala@bala.bg",
		Msg:      spf(`Gotowo: <a href://"site.com/%s">Ime na kniga</a>`, downloadFileName),
	}
	os.MkdirAll(dir, 0755)
	// TODO: add most of these to DefaultConfig
	cmdArgs := []string{
		spf("За %s<%s>", pdfMsg.Name, pdfMsg.Email),
		"font:Pliska-Italic, points:9, off:0 20, sc:1.0 abs, pos:bc, rot:0",
		sourceFileName,
		downloadFileName,
		passw,
	}
	c.Echo().Logger.Debugf("pdfcpu.sh %v", cmdArgs)
	cmd := exec.Command("./bin/pdfcpu_stamp_encrypt.sh", cmdArgs...)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusCreated, pdfMsg)
}

// pdfcpuMessage is produced by ppdfcpu()
type pdfcpuMessage struct {
	Download string `json:"download" xml:"download" form:"download" query:"download"`
	Passwd   string `json:"passwd" xml:"passwd" form:"passwd" query:"passwd"`
	Name     string `json:"name" xml:"name" form:"name" query:"name"`
	Email    string `json:"email" xml:"email" form:"email" query:"email"`
	Msg      string `json:"msg" xml:"msg" form:"msg" query:"msg"`
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
	c.Logger().Debugf("in hello")
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
// strings. This map is used in loadRoutes() to match HTTP handlerFuncs by name.
var handlerFuncs = map[string]echo.HandlerFunc{
	"hello":       hello,
	"ppdfcpu":     ppdfcpu,
	"ppdfcpuForm": ppdfcpuForm,
}

// This map is for the same purpuse as above but for one or more middleware
// functions for the corresponding HandlerFunc.
var middlewareFuncs = map[string]echo.MiddlewareFunc{}
