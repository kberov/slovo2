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
	log.Debugf("in straniciExecute")
	args := new(model.StranicaArgs)
	if err := c.Bind(args); err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	page := new(model.Stranici)
	if err := page.FindForDisplay(args); err != nil {
		log.Errorf("page: %#v; error:%s", page, err)
		return err
	}
	stash := Map{
		"title":      page.Title,
		"page.Alias": page.Alias,
		"page.ID":    spf("%d", page.ID),
		"page.Body":  page.Body,
	}
	stash["main_menu"] = prepareMenu(c, args, stash)
	return c.Render(200, "stranici/execute", stash)
}

func prepareMenu(c echo.Context, args *model.StranicaArgs, stash Map) gledki.TagFunc {
	return gledki.TagFunc(func(w io.Writer, tag string) (int, error) {
		items, err := model.SelectMenuItems(args)
		if err != nil {
			c.Logger().Error(err.Error())
			return w.Write([]byte("error retrieving items... see log for details"))
		}
		html := bytes.NewBuffer([]byte(""))
		for _, p := range items {
			class := ""
			if p.Alias == stash["page.Alias"] {
				class = `class="active"`
			}
			html.WriteString(spf(`<a %s href="/%s.%s.html">%s</a>`, class, p.Alias, p.Language, p.Title))
		}
		return w.Write(html.Bytes())
	})
}
