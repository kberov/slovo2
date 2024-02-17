package slovo

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

const cached = `cached`

// Cache pages so when Apache finds a ready page, slovo is not invoked at all.
func cachePages(ec echo.Context, reqBody, resBody []byte) {
	c := ec.(*Context)
	// If the user is not Guest or the file is not html, do not cache!
	if !Cfg.CachePages && (c.StraniciArgs.Format != StraniciFormat || c.StraniciArgs.UserID != GuestID) {
		return
	}

	path := pathToFile(c)
	// c.Logger().Debugf("filePath:%s", path)
	domain, _ := strings.CutPrefix(hostName(c), `dev`)
	domain, _ = strings.CutPrefix(hostName(c), `www`)
	fullPath := filepath.Join(BinDir(), `domove`, domain, `public`, cached, path)
	// c.Logger().Debugf("fullPath: %s", fullPath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0700); err != nil {
		c.Logger().Panic(err)
	}
	if err := os.WriteFile(fullPath, resBody, 0600); err != nil {
		c.Logger().Panic(err)
	}
}

func pathToFile(c echo.Context) string {
	return c.Request().RequestURI
}
