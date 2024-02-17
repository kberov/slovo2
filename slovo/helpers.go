package slovo

import (
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/idna"
)

// helper functions for the slovo package

// hostName extracts punycode encoded domain name and returns it. For decoded
// unicode domain name use iHostName(c)
func hostName(c echo.Context) (host string) {
	host = c.Request().Host
	if !strings.Contains(host, ":") {
		return
	}
	host, _, err := net.SplitHostPort(host)
	if err != nil {
		c.Logger().Errorf("could not parse host from %s; err:%s", c.Request().Host, err)
	}
	return
}

// iHostName(c) returns converted to unicode domain name for displaying it on
// pages.
func iHostName(c echo.Context) (host string) {
	host, err := idna.New().ToUnicode(hostName(c))
	if err != nil {
		c.Logger().Errorf("could not parse international host from %s; err:%s", c.Request().Host, err)
	}
	return
}

// Allow only valid values 0,1,2
func publishedStatus(c echo.Context) uint8 {
	preview := c.QueryParam("preview")
	if preview != "" {
		return uint8(1)
	}
	return uint8(2)
}

var binDir string

// BinDir returns the directory where the binary resides. We assume that this
// is the root of the project.
func BinDir() string {
	if len(binDir) > 0 {
		return binDir
	}
	exe, err := os.Executable()
	if err != nil {
		panic(err)
	}
	binDir = filepath.Dir(exe)
	return binDir
}
