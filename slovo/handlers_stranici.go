package slovo

import (
	"bytes"
	"io"
	"net/http"

	"github.com/kberov/gledki"
	"github.com/kberov/slovo2/model"
	"github.com/labstack/echo/v4"
)

func straniciExecute(c echo.Context) error {
	log := c.Logger()
	args := new(model.StranicaArgs)
	if err := c.Bind(args); err != nil {
		return c.String(http.StatusBadRequest, "Грешна заявка!")
	}
	page := new(model.Stranici)
	if err := page.FindForDisplay(args); err != nil {
		log.Errorf("page: %#v; error:%s", page, err)
		return err
	}
	// log.Debugf("Stranica: %#v", page)
	stash := Map{
		"lang":       page.Language,
		"title":      page.Title,
		"page.Alias": page.Alias,
		"page.ID":    spf("%d", page.ID),
	}
	stash["pageBody"] = pageBody(c, page, stash)
	stash["mainMenu"] = mainMenu(c, args, stash)

	return c.Render(200, page.TemplatePath("stranici/execute"), stash)
}

/*
mainMenu returns a gledki.TagFunc which prepares and returns the HTML for
the tag `mainMenu` in the template.
*/
func mainMenu(c echo.Context, args *model.StranicaArgs, stash Map) gledki.TagFunc {
	return func(w io.Writer, tag string) (int, error) {
		items, err := model.SelectMenuItems(args)
		if err != nil {
			c.Logger().Errorf(`error from model.SelectMenuItems: %w`, err)
			return w.Write([]byte("error retrieving items... see log for details"))
		}
		html := bytes.NewBuffer([]byte(""))
		for _, p := range items {
			class := ""
			if p.Alias == stash["page.Alias"] {
				class = `class="active" `
			}
			html.WriteString(spf(`<a %shref="/%s.%s.html">%s</a>`, class, p.Alias, p.Language, p.Title))
		}
		return w.Write(html.Bytes())
	}
}

/*
pageBody returns a gledki.TagFunc which prepares and returns the HTML for
the tag `pageBody` in the template.
*/
func pageBody(c echo.Context, page *model.Stranici, stash Map) gledki.TagFunc {
	// TODO: Implement custom logic depending on the page.Template which has to be filled in.
	// It has to work somehow automatically. We should not have to write new
	// code if new template is added in the site, or maybe have a limited set
	// of templates which can be chosen from a select<options> dropdown in the
	// control panel.
	return func(w io.Writer, tag string) (int, error) {
		switch page.Template.String {
		case `stranici/templates/dom`:
			return w.Write([]byte(page.Body))
		default:
			return w.Write([]byte(page.Body))
		}
	}
}
