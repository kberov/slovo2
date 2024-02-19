package slovo

import (
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kberov/gledki"
	"github.com/kberov/slovo2/model"
	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
	StraniciArgs  *model.StraniciArgs
	canonicalPath string
}

func (c *Context) DB() *sqlx.DB {
	return model.DB()
}

// BindArgs prepares common arguments for `stranici` and `celini`. It is
// idempotent. If invoked multiple times, returns the same prepared
// [model.StraniciArgs].
func (c *Context) BindArgs() (*model.StraniciArgs, error) {
	if c.StraniciArgs.UserID > 0 {
		return c.StraniciArgs, nil
	}
	err := c.Bind(c.StraniciArgs)
	return c.StraniciArgs, err
}

// CanonicalPath returns the canonical URL for the current page.
func (c *Context) CanonicalPath() string {
	if c.canonicalPath != "" {
		return c.canonicalPath
	}
	var path strings.Builder
	path.WriteByte('/')
	path.WriteString(c.StraniciArgs.Alias)
	if c.StraniciArgs.Celina != "" {
		path.WriteByte('/')
		path.WriteString(c.StraniciArgs.Celina)
	}
	path.WriteByte('.')
	path.WriteString(c.StraniciArgs.Lang)
	path.WriteByte('.')
	path.WriteString(c.StraniciArgs.Format)
	c.canonicalPath = path.String()
	return c.canonicalPath
}

func (c *Context) prepareDefaultStash() {
	c.Echo().Renderer.(*EchoRenderer).MergeStash(gledki.Stash{
		"Date":      time.Now().Format(time.RFC1123),
		"canonical": "https://" + c.StraniciArgs.Domain + c.CanonicalPath(),
		"generator": "slovo2",
		"version":   VERSION,
		"codename":  CODENAME,
	})
}

/*
SlovoContext is a middleware function which instantiates slovo's custom context
and executes some tascs common to all pages in the site. These are:
  - [Context.BindArgs]
  - renders cached pages
  - prepares some default items in [gledki.Stash]
*/
func SlovoContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Logger().Debugf("in SlovoContext")
		sc := &Context{Context: c, StraniciArgs: new(model.StraniciArgs)}
		if _, err := sc.BindArgs(); err != nil {
			return err
		}
		if canCachePage(sc) {
			if err := tryHandleCachedPage(sc); err == nil {
				return nil
			}
		}
		sc.prepareDefaultStash()
		return next(sc)
	}
}
