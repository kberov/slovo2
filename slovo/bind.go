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

func (b *Binder) Bind(args any, c echo.Context) (err error) {
	/* Using default binder */
	if err = b.DefaultBinder.Bind(args, c); err != nil {
		return err
	}

	// Define our custom implementation here
	// TODO: See Echo.DefaultBinder for implementation details and follow it's
	// pattern if better for the case.
	a := args.(*model.StranicaArgs)
	// TODO implement authentication and see if we need the whole user somewhere.
	// user := new(model.Users)
	// model.GetByID(user, 2)
	// a.UserID = user.ID
	a.UserID = 2
	a.Pub = publishedStatus(c)
	a.Domain = hostName(c)
	a.Now = time.Now().Unix()
	return
}
