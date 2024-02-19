package slovo

import (
	"net/http"
	"regexp"
	"time"

	"github.com/kberov/slovo2/model"
	"github.com/labstack/echo/v4"
)

func celiniExecute(ec echo.Context) error {
	c := ec.(*Context)
	log := c.Logger()
	log.Debugf("in celiniExecute")
	cel := new(model.Celini)
	if err := cel.FindForDisplay(*c.StraniciArgs); err != nil {
		log.Errorf("celina: %#v; error:%w; ErrType: %T; args: %#v", cel, err, err, c.StraniciArgs)
		return handleNotFound(c, err)
	}
	return c.Render(http.StatusOK, cel.TemplatePath("celini/note"), buildCeliniStash(c, cel))
}

func buildCeliniStash(c *Context, cel *model.Celini) Stash {
	args := c.StraniciArgs
	user := new(model.Users)
	model.GetByID(user, cel.UserID)
	created := time.Unix(int64(cel.CreatedAt), 0)
	tstmp := time.Unix(int64(cel.Tstamp), 0)
	stash := Stash{
		"lang":       cel.Language,
		"title":      cel.Title,
		"page.Alias": args.Alias,
		"ogType":     "article",
		"cel.ID":     spf("%d", cel.ID),
		"UserNames":  user.FirstName + ` ` + user.LastName,
		"CreatedAt":  created.Format(time.DateOnly),
		"Tstamp":     tstmp.Format(time.DateOnly),
	}
	stash["mainMenu"] = mainMenu(c, args, stash)
	stash["celBody"] = celBody(c, cel, stash)

	return stash
}

/*
celBody returns a gledki.TagFunc which prepares and returns the HTML for
the tag `celBody` in the template.
*/
func celBody(c *Context, cel *model.Celini, stash Stash) string {
	// prepare different values for the stash depnding on the DataType
	switch cel.DataType {
	case model.Book:
		return cel.Body
	default:
		stash["ogImage"] = ogImage(c, cel.Body)
		return cel.Body
	}
}

var reOgImage = regexp.MustCompile(`(?i:<img.+?src="([^"]+\.(?:png|jpe?g|webp)))"`)

/*
ogImage finds the first image tag in the celBody string and returns the value
of its src attribute. If not found, returns an empty string
*/
func ogImage(c *Context, celBody string) string {
	match := reOgImage.FindStringSubmatch(celBody)
	if len(match) > 0 {
		return match[1]
	}
	return ""
}
