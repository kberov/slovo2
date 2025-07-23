package slovo

import (
	"errors"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

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
func publishedStatus(c echo.Context) int {
	preview := c.QueryParam("preview")
	if preview != "" {
		return 1
	}
	return 2
}

func FileIsReadable(path string) bool {
	finfo, err := os.Stat(path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false
	}
	if finfo.Mode().IsRegular() && finfo.Mode().Perm()&0400 == 0400 {
		return true
	}
	return false
}

var reHTML = regexp.MustCompile(`<[^>]+>`)

func stripHTML(text string) string {
	return reHTML.ReplaceAllString(text, "")
}

/*
substring extracts a substring out of `expr` and returns it. First character
is at offset zero. If LENGTH is 0, returns everything through the end of the
string. String is a string of runes.
*/

func substring(expr string, offset uint, length uint) string {
	characters := utf8.RuneCountInString(expr)
	if length == 0 {
		return expr
	}
	if uint(characters) < offset+length {
		return expr
	}
	return string([]rune(expr)[offset:length])
}

/*
substringWithTail does the same as substring, but adds a tail string in case
the input string was longer than the output string.
*/
func substringWithTail(expr string, offset uint, length uint, tail string) string {
	if utf8.RuneCountInString(expr) > int(length) {
		return substring(expr, offset, length) + tail
	}
	return expr
}

// domainName return the current domain name without common prefixes like
// dev,www,qa etc, as listed in Cfg.DomovePrefixes.
func domainName(c echo.Context) string {
	domainName, ok := c.Get(`domainName`).(string)
	if ok {
		return domainName
	}
	for _, prefix := range Cfg.DomovePrefixes {
		domainName, isCut := strings.CutPrefix(hostName(c), prefix)
		if isCut {
			c.Set(`domainName`, domainName)
			return domainName
		}
	}
	return domainName
}

var binDir string

// HomeDir returns the slovo2 installation directory. This is the directory
// where we have the `domove` directory.
func HomeDir() string {
	if binDir != "" {
		return binDir
	}
	cwd, _ := os.Getwd()
	dir := `domove`
	path := filepath.Join(cwd, dir)
	for {
		_, err := os.Stat(path)
		if err != nil && errors.Is(err, os.ErrNotExist) {
			path = filepath.Dir(path)
			path = filepath.Join(filepath.Dir(path), dir)
			continue
		}

		if path == "" {
			panic("HomeDir could not be found")
		}
		binDir = filepath.Dir(path)
		return binDir
	}
}
