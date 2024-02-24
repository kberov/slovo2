package slovo

import (
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	m "github.com/kberov/slovo2/model"
	"github.com/labstack/echo/v4"
)

func straniciExecute(ec echo.Context) error {
	c := ec.(*Context)
	log := c.Logger()
	page := new(m.Stranici)
	if err := page.FindForDisplay(*c.StraniciArgs); err != nil {
		log.Errorf("%v; ErrType: %T; args: %#v", err, err, c.StraniciArgs)
		return handleNotFound(c, err)
	}
	return c.Render(http.StatusOK, page.TemplatePath("stranici/execute"), buildStraniciStash(c, page))
}

// buildStraniciStash adds all the needed tags to be replaced in template with their
// values. Returns the prepared stash - a map["string"]any.
func buildStraniciStash(c *Context, page *m.Stranici) Stash {
	args := c.StraniciArgs
	stash := Stash{
		"lang":       page.Language,
		"title":      page.Title,
		"page.Alias": page.Alias,
		"page.ID":    spf("%d", page.ID),
		"ogType":     "website",
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
		stash["categoryPages"] = categoryPages(c, stash)
	// other cases maybe
	// case`stranici/other/special/view`
	default:
		stash["categoryCelini"] = categoryCelini(c, stash)
	}
	return stash
}

/*
mainMenu returns a gledki.TagFunc which prepares and returns the HTML for
the tag `mainMenu` in the template.
*/
func mainMenu(c echo.Context, args *m.StraniciArgs, stash Stash) string {
	var html strings.Builder
	for _, p := range m.SelectMenuItems(*args) {
		class := ""
		if p.Alias == stash["page.Alias"] {
			class = `class="active" `
		}
		html.WriteString(spf(`<a %shref="/%s.%s.html">%s</a>`, class, p.Alias, p.Language, p.Title))
	}
	return html.String()
}

// categoryPages displays the list of pages in the home page.
func categoryPages(c *Context, stash Stash) string {
	t, _ := c.Echo().Renderer.(*EchoRenderer)

	// File does not have directives in it self, so only LoadFile() is
	// enough. No need to Compile().
	partial := t.MustLoadFile(`stranici/_dom_item`)
	var view strings.Builder
	for _, page := range m.ListStranici(*c.StraniciArgs) {
		stash := Stash{
			"id":    spf("%d", page.ID),
			"title": page.Title,
			"lang":  page.Language,
			"alias": page.Alias,
			"body":  substringWithTail(stripHTML(page.Body), 0, 220, ` …`),
		}
		view.WriteString(t.FtExecStringStd(partial, stash))
	}
	return view.String()
}

// categoryCelini displays the list of celini in the respective category page.
func categoryCelini(c *Context, stash Stash) string {
	t, _ := c.Echo().Renderer.(*EchoRenderer)

	partialT := t.MustLoadFile("stranici/_cel_item")
	var view strings.Builder
	celini := m.ListCelini(*c.StraniciArgs)
	for _, cel := range celini {
		hash := Stash{
			"id":        spf("%d", cel.ID),
			"title":     substringWithTail(cel.Title, 0, 24, `…`),
			"fullTitle": cel.Title,
			"body":      substring(stripHTML(cel.Body), 0, 200) + " …",
			"alias":     cel.Alias,
			"strAlias":  c.StraniciArgs.Alias,
			"lang":      cel.Language,
		}
		view.WriteString(t.FtExecStringStd(partialT, hash))
	}
	stash[`categoryCeliniPager`] = categoryCeliniPager(c, len(celini))
	return view.String()
}

// categoryCeliniPager displays `<`:previous and `>`:next links under the list
// of celini.
func categoryCeliniPager(c *Context, celiniNum int) string {
	args := c.StraniciArgs
	if celiniNum < args.Limit && args.Offset == 0 {
		return ``
	}
	// TODO: move the HTML to a partial template and use
	// t.FtExecStringStd(partial, stash)! See categoryPages.
	pagertmpl := `<div class="card col-12 pager">`
	// link to previous
	if args.Offset > 0 {
		offset := args.Offset - args.Limit
		if offset <= 0 {
			pagertmpl += spf(`<a title="първи %[4]d" href="/%[1]s.%s.%s">⮈</a>`,
				args.Alias, args.Lang, args.Format, args.Limit)
		} else {
			pagertmpl += spf(`<a title="предишни %[4]d" href="/%[1]s.%s.%s?limit=%d&offset=%d">⮈</a>`,
				args.Alias, args.Lang, args.Format, args.Limit, offset)
		}
		if celiniNum == args.Limit {
			pagertmpl += `&nbsp;&nbsp;`
		}
	}
	// link to next
	if celiniNum == args.Limit {
		pagertmpl += spf(`<a title="следващи %[4]d" href="/%[1]s.%s.%s?limit=%d&offset=%d">⮊</a>`,
			args.Alias, args.Lang, args.Format, args.Limit, (args.Offset + args.Limit))
	}
	return pagertmpl + `</div>`
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

/*
substringWithTail does the same as substring, but adds a tail string in case
the input string was longer than the output string.
*/
func substringWithTail(expr string, offset uint, length uint, tail string) string {
	if utf8.RuneCountInString(expr) > int(length) {
		return substring(expr, offset, length) + tail
	}
	return expr
}
