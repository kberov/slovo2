package slovo

import (
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
	// proper SQL queries.
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
		a.Domain = hostName(c)
		a.Now = time.Now().Unix()
		// By default the main box is displayed as the main content on the
		// rendered page.
		a.Box = model.MainBox
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
