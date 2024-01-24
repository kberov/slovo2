package slovo

import (
	"net"
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