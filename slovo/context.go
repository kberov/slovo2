package slovo

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/kberov/gledki"
	m "github.com/kberov/slovo2/model"
	"github.com/labstack/echo/v4"
)

/*
Context is our custom context. It extends [echo.Context]. To use it in our
handlers, we have to cast echo.Context to slovo.Context. Here is a full example.

	func celiniExecute(ec echo.Context) error {
		c := ec.(*Context)
		log := c.Logger()
		cel := new(model.Celini)
		if err := cel.FindForDisplay(*c.StraniciArgs); err != nil {
			log.Errorf("celina: %#v; error:%w; ErrType: %T; args: %#v", cel, err, err, c.StraniciArgs)
			return handleNotFound(c, err)
		}
		return c.Render(http.StatusOK, cel.TemplatePath("celini/note"), buildCeliniStash(c, cel))
	}
*/
type Context struct {
	echo.Context
	// StraniciArgs contains a pointer to [m.StraniciArgs].
	StraniciArgs  *m.StraniciArgs
	canonicalPath string
	// DomainRoot is the root folder for static content folders and `templates`
	// folder of the domain for this HTTP request.
	DomainRoot string
	// Domain contains a pointer to a m.Domove record - the current domain.
	Domain *m.Domove
}

func (c *Context) DB() *sqlx.DB {
	return m.DB()
}

// BindArgs prepares common arguments for `stranici` and `celini`. It is
// idempotent. If invoked multiple times within a request, it returns the same
// prepared [model.StraniciArgs].
func (c *Context) BindArgs() (*m.StraniciArgs, error) {
	if c.StraniciArgs.UserID > 0 {
		return c.StraniciArgs, nil
	}
	err := c.Bind(c.StraniciArgs)

	// Make sure we use only the domain name without any prefix.
	дом := new(m.Domove)
	дом.GetByName(c.StraniciArgs.Domain)
	c.StraniciArgs.Domain = дом.Domain
	c.DomainRoot = filepath.Join(Cfg.DomoveRoot, c.StraniciArgs.Domain)
	//c.Logger().Debugf("Domain:%#v", dom)
	c.Domain = дом
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

func (c *Context) switchToDomainTemplates() {
	// Check if the current domain has set its own templates at all.
	dom := c.Domain
	if dom.Templates == "" {
		return
	}
	domainTemlatesRoot := filepath.Join(Cfg.DomoveRoot, dom.Domain, "templates")
	domainThemeRoot := filepath.Join(domainTemlatesRoot, dom.Templates)
	gl := c.Echo().Renderer.(*EchoRenderer)
	c.Logger().Debugf("Template roots: %#v", gl.Roots)
	// Check if templates are already at first and second place in Roots.
	// Prepended during the previous request which happened to be to the same
	// domain.
	if len(gl.Roots) > 2 && gl.Roots[0] == domainThemeRoot && gl.Roots[1] == domainTemlatesRoot {
		return
	}
	// Check if there is already a template root in domove.
	if !dirExists(domainThemeRoot) {
		c.Logger().Debugf("Domain %s has set a theme '%s', but template root '%s' "+
			" does not exist. You may want to create it.", dom.Domain, dom.Templates, domainTemlatesRoot)
		return
	}
	// To be safe, remove any root which is under Cfg.DomoveRoot
	gl.Roots = slices.DeleteFunc(gl.Roots, func(path string) bool {
		return strings.Contains(path, Cfg.DomoveRoot)
	})
	// Now prepend the new roots to be searched for templates
	gl.Roots = slices.Insert(gl.Roots, 0, domainThemeRoot, domainTemlatesRoot)
	c.Logger().Debugf("Template roots: %#v", gl.Roots)
}

/*
SlovoContext is a middleware function which instantiates slovo's custom context
and executes some tasks common to all pages in the site. These are:
  - [Context.BindArgs]
  - renders (spits out) cached pages
  - sets the domain root c.DomainRoot after the current domain for this HTTP request.
  - modifies gledki.Gledki.TemplatesRoots so the first path is always defined
    by the current domain if the `domove` table record contains some value in its
    templates column. This way every domain can have it's own templates.
  - prepares some default items in [gledki.Stash]
*/
func SlovoContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Logger().Debugf("in SlovoContext")
		sc := &Context{Context: c, StraniciArgs: new(m.StraniciArgs)}
		if _, err := sc.BindArgs(); err != nil {
			return err
		}
		if err := tryHandleCachedPage(sc); err == nil {
			return nil
		}
		sc.switchToDomainTemplates()
		sc.prepareDefaultStash()
		return next(sc)
	}
}

// copied from gledki
func dirExists(path string) bool {
	finfo, err := os.Stat(path)
	if err != nil && errors.Is(err, os.ErrNotExist) || !finfo.IsDir() {
		return false
	}
	return true
}
