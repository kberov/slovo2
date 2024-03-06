package slovo

import (
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/labstack/echo/v4"
)

const cached = `cached`

// Cache pages so when Apache finds a ready page, slovo is not invoked at all.
// Invoked by echo middleware.BodyDump().
func cachePages(ec echo.Context, reqBody, resBody []byte) {
	c := ec.(*Context)
	if !canCachePage(c) {
		return
	}
	// If page is already cached, do not overwrite the file with the just
	// extracted from it content.
	path := c.CanonicalPath()
	fullPath := filepath.Join(Cfg.DomoveRoot, c.StraniciArgs.Domain, `public`, cached, path)
	if FileIsReadable(fullPath) {
		return
	}
	// c.Logger().Debugf("in cachePages filePath:%s", path)
	// c.Logger().Debugf("fullPath: %s", fullPath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		c.Logger().Panic(err)
	}
	if err := os.WriteFile(fullPath, resBody, 0644); err != nil {
		c.Logger().Panic(err)
	}
}

// canCachePage says if the current page can be cached. Uses slovo.Context
func canCachePage(c *Context) bool {
	// If caching is not enabled, do not cache!
	if !Cfg.CachePages {
		return false
	}
	// Cache only GET requests
	if c.Request().Method != http.MethodGet {
		return false
	}
	// Cache only requests without parameters.
	if c.QueryString() != "" {
		return false
	}
	// If the user is not Guest or the file is not html, do not cache!
	if c.StraniciArgs.Format != format || c.StraniciArgs.UserID != Cfg.GuestID {
		return false
	}
	return true
}

// PreferDomainStaticFiles serves static files from domain specific
// directories if these files exist.
func PreferDomainStaticFiles(next echo.HandlerFunc) echo.HandlerFunc {
	reStaticFile := regexp.MustCompile(Cfg.DomoveStaticFiles)
	return func(c echo.Context) error {
		path := c.Request().URL.Path
		c.Logger().Debugf("request file: %s", path)
		if reStaticFile.MatchString(path) {
			domain := domainName(c)
			file := filepath.Join(Cfg.DomoveRoot, domain, `public`, path)
			c.Logger().Debugf("filepath:%s", file)
			if FileIsReadable(file) {
				return c.File(file)
			}
		}
		return next(c)
	}
}
