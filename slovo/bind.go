package slovo

import (
	"strings"
	"time"

	"github.com/kberov/slovo2/model"
	"github.com/labstack/echo/v4"
)

// Binder embeds echo.DefaultBinder
// TODO: Read whole https://go101.org/article/type-embedding.html
type Binder struct {
	*echo.DefaultBinder
}

/*
Bind binds some variables into a structure to be passed to queries for
Stranici and Celini.
*/
func (b *Binder) Bind(args any, c echo.Context) (err error) {
	// Here we handle untagged fields - those which values cannot be simply got
	// from any of the supported by [echo] tags. But we need them to make
	// proper SQL queries. We also set default values for some potentially
	// tagget fiedls like Limit and Offset.
	switch t := args.(type) {
	case *model.StraniciArgs:
		a := t
		// TODO implement authentication and see if we need the whole user somewhere.
		// user := new(model.Users)
		// Default user - guest
		// model.GetByID(user, 2)
		// a.UserID = user.ID
		a.UserID = 2
		a.Pub = publishedStatus(c)
		a.Domain, _ = strings.CutPrefix(hostName(c), `dev.`)
		a.Domain, _ = strings.CutPrefix(a.Domain, `www.`)
		a.Now = time.Now().Unix()
		// By default the main box is displayed as the main content on the
		// rendered page.
		a.Box = model.MainBox
		a.Limit = 100
		a.Offset = 0
	//case *model.SomeOtherArgs:
	//	a := t
	// etc
	default:
		c.Logger().Warnf("Unknown type: %T", args)
	}

	/* Using default binder */
	if err = b.DefaultBinder.Bind(args, c); err != echo.ErrUnsupportedMediaType {
		return err
	}

	return
}
