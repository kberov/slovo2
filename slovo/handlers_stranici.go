package slovo

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/kberov/gledki"
	"github.com/kberov/slovo2/model"
	"github.com/labstack/echo/v4"
)

func straniciExecute(c echo.Context) error {
	args := new(model.StraniciArgs)
	if err := c.Bind(args); err != nil {
		return c.String(http.StatusBadRequest, "Грешна заявка!")
	}

	page := new(model.Stranici)
	if err := page.FindForDisplay(*args); err != nil {
		c.Logger().Errorf("page: %#v; error:%w; ErrType: %T; args: %#v", page, err, err, args)
		return handleNotFound(c, args, err)
	}
	stash := Map{
		"lang":       page.Language,
		"title":      page.Title,
		"page.Alias": page.Alias,
		"page.ID":    spf("%d", page.ID),
	}
	stash["pageBody"] = pageBody(c, page, stash)
	stash["mainMenu"] = mainMenu(c, args, stash)
	stash["categoryPages"] = categoryPages(c, *args, stash)

	return c.Render(http.StatusOK, page.TemplatePath("stranici/execute"), stash)
}

func handleNotFound(c echo.Context, args *model.StraniciArgs, err error) error {
	stash := Map{"lang": args.Lang, "title": "Няма такава страница!"}
	if strings.Contains(err.Error(), `no rows`) {
		stash["mainMenu"] = mainMenu(c, args, stash)
		return c.Render(http.StatusNotFound, `not_found`, stash)
	}
	return err
}

/*
mainMenu returns a gledki.TagFunc which prepares and returns the HTML for
the tag `mainMenu` in the template.
*/
func mainMenu(c echo.Context, args *model.StraniciArgs, stash Map) gledki.TagFunc {
	return func(w io.Writer, tag string) (int, error) {
		items, err := model.SelectMenuItems(*args)
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

// categoryPages displays the list of pages in the home page.
func categoryPages(c echo.Context, args model.StraniciArgs, stash Map) gledki.TagFunc {
	t, ok := c.Echo().Renderer.(*EchoRenderer)
	if !ok {
		err := errors.New(spf("slovo2 works only with the `gledki` template engine. This is %T", c.Echo().Renderer))
		c.Logger().Error(err)
		return func(w io.Writer, tag string) (int, error) {
			return 0, err
		}
	}

	return func(w io.Writer, tag string) (int, error) {
		templatePath := `stranici/_dom_item`
		partialTemplate, err := t.Compile(templatePath)
		if err != nil {
			return 0, err
		}
		pagesBB := bytes.NewBuffer([]byte(""))
		childrenPages, err := model.ListStranici(args)
		if err != nil {
			c.Logger().Error(err)
			return w.Write(pagesBB.Bytes())
		}
		for _, page := range childrenPages {
			stash := Map{
				"id":    spf("%d", page.ID),
				"title": page.Title,
				"lang":  page.Language,
				"alias": page.Alias,
				"body":  substring(stripHTML(page.Body), 0, 220),
			}
			if _, err := t.FtExecStd(partialTemplate, pagesBB, stash); err != nil {
				return 0, fmt.Errorf("Problem rendering partial template %s TagFunc: %s", partialTemplate, err.Error())
			}
		}
		return w.Write(pagesBB.Bytes())
	}
}

var reHTML = regexp.MustCompile(`<[^>]+>`)

func stripHTML(text string) string {
	return reHTML.ReplaceAllString(text, "")
}

/*
substring extracts a substring out of `expr` and returns it. First character
is at offset zero. If LENGTH is 0, returns everything through the end of the
string. String is a string of runes.
*/

func substring(expr string, offset uint, length uint) string {
	characters := utf8.RuneCountInString(expr)
	if length == 0 {
		return expr
	}
	if uint(characters) < offset+length {
		return expr
	}
	return string([]rune(expr)[offset:length])
}
