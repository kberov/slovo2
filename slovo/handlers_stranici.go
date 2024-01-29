package slovo

import (
	"github.com/kberov/slovo2/model"
	"github.com/labstack/echo/v4"
)

func straniciExecute(c echo.Context) error {
	log := c.Logger()
	log.Debugf("in straniciExecute")
	lang := c.Param("lang")
	/*if lang == "" && len(c.Request().Header["Accept-Language"]) > 0 {
		lang = c.Request().Header["Accept-Language"][0]
	}
	*/
	pageAlias := c.Param("stranica")
	path := c.Request().URL.Path
	data := spf("Params: pagealias: %s, lang:%s; Path:%s; PathInfo:%s; QUERY_STRING:%s",
		pageAlias, lang, c.Path(), path, c.QueryString())
	c.Logger().Debug(data)
	// TODO: Implement sessions for users, using NOT Cookies, but something
	// else - header+localStorage, JWT... who knows
	user := new(model.Users)
	model.GetByID(user, 2)
	page := new(model.Stranici)
	domain := hostName(c)
	preview := publishedStatus(c)
	if err := page.FindForDisplay(pageAlias, user, preview, domain, lang); err != nil {
		log.Errorf("page: %#v; error:%s", page, err)
		return err
	}
	//now we have in page a structure
	log.Debugf("page: %#v; unicode domain: %s", page, iHostName(c))

	/* not needed as we select the content with one JOIN sql statement in page.FindForDisplay
	// get the appropriate celina for this page
	cel := new(model.Celini)
	if err := cel.FindPageCelinaForDisplay(page, user, preview, lang[:2]); err != nil {
		log.Errorf("celina: %#v; error:%s", cel, err)
		return err
	}
	log.Debugf("celina: %#v;", cel)
	*/
	return c.Render(200, "stranici/execute",
		Map{
			"title":      "Страница еди коя си!",
			"greeting":   "Добре дошли! на страница " + data,
			"page.Alias": page.Alias,
			"page.ID":    spf("%d", page.ID),
		},
	)
}
