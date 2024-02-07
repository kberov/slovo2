package slovo

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/kberov/gledki"
	"github.com/kberov/slovo2/model"
	"github.com/labstack/echo/v4"
)

const (
	wrongRendererMsg = `slovo2 works only with the "gledki" template engine. This is %T`
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
	return c.Render(http.StatusOK, page.TemplatePath("stranici/execute"), buildStraniciStash(c, page, args))
}

func handleNotFound(c echo.Context, args *model.StraniciArgs, err error) error {
	// TODO: I18N & L10N
	stash := Stash{"lang": args.Lang, "title": "Няма такава страница!"}
	if strings.Contains(err.Error(), `no rows`) {
		stash["mainMenu"] = mainMenu(c, args, stash)
		return c.Render(http.StatusNotFound, `not_found`, stash)
	}
	return err
}

// buildStraniciStash adds all the needed tags to be replaced in template wit their
// values. Returns the prepared stash - a map["string"]any.
func buildStraniciStash(c echo.Context, page *model.Stranici, args *model.StraniciArgs) Stash {
	stash := Stash{
		"lang":       page.Language,
		"title":      page.Title,
		"page.Alias": page.Alias,
		"page.ID":    spf("%d", page.ID),
	}
	stash["mainMenu"] = mainMenu(c, args, stash)
	stash["pageBody"] = page.Body
	/*
	   TODO: when needed, implement custom logic depending on the page.Template
	   which has to be filled in.  It has to work somehow automatically. We
	   should not have to write new code if new template is added in the site,
	   or maybe have a limited set of templates which can be chosen from a
	   select<options> dropdown in the control panel and have some mechanism to
	   bind code to templates. We actually already have it with the TagFunc
	   concept from fasttemplate.
	*/
	switch page.Template.String {
	case `stranici/templates/dom`:
		stash["categoryPages"] = categoryPages(c, *args, stash)
	// other cases maybe
	// case`stranici/other/special/view`
	default:
		stash["categoryCelini"] = categoryCelini(c, *args, stash)
	}
	return stash
}

/*
mainMenu returns a gledki.TagFunc which prepares and returns the HTML for
the tag `mainMenu` in the template.
*/
func mainMenu(c echo.Context, args *model.StraniciArgs, stash Stash) gledki.TagFunc {
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

// categoryPages displays the list of pages in the home page.
func categoryPages(c echo.Context, args model.StraniciArgs, stash Stash) string {
	t, ok := c.Echo().Renderer.(*EchoRenderer)
	if !ok {
		err := errors.New(spf(wrongRendererMsg, c.Echo().Renderer))
		c.Logger().Error(err)
		return ""
	}

	// File does not have directives in it self, so only LoadFile() is
	// enough. No need to Compile().
	templatePath := `stranici/_dom_item`
	partialTemplate, err := t.LoadFile(templatePath)
	if err != nil {
		c.Logger().Error(err)
		return ""
	}
	childrenPages, err := model.ListStranici(args)
	if err != nil {
		c.Logger().Error(err)
		return ""
	}
	var view strings.Builder
	for _, page := range childrenPages {
		stash := Stash{
			"id":    spf("%d", page.ID),
			"title": page.Title,
			"lang":  page.Language,
			"alias": page.Alias,
			"body":  substring(stripHTML(page.Body), 0, 220),
		}
		view.WriteString(t.FtExecStringStd(partialTemplate, stash))
	}
	return view.String()
}

// categoryCelini displays the list of celini in the respective category page.
func categoryCelini(c echo.Context, args model.StraniciArgs, stash Stash) string {

	t, ok := c.Echo().Renderer.(*EchoRenderer)
	if !ok {
		err := errors.New(spf(wrongRendererMsg, c.Echo().Renderer))
		c.Logger().Error(err)
		return ""
	}
	partialT, err := t.LoadFile("stranici/_cel_item")
	if err != nil {
		c.Logger().Error(err)
		return ""
	}
	celini, err := model.ListCelini(args)
	if err != nil {
		c.Logger().Error(err)
		return ""
	}
	var view strings.Builder
	for _, cel := range celini {
		title := ""
		if utf8.RuneCountInString(cel.Title) > 26 {
			title = substring(cel.Title, 0, 26) + "…"
		} else {
			title = cel.Title
		}
		hash := Stash{
			"id":        spf("%d", cel.ID),
			"title":     title,
			"fullTitle": cel.Title,
			"body":      substring(stripHTML(cel.Body), 0, 200) + "…",
			"alias":     cel.Alias,
			"strAlias":  args.Alias,
			"lang":      cel.Language,
		}
		view.WriteString(t.FtExecStringStd(partialT, hash))
	}
	return view.String()
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
