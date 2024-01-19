package slovo

import "github.com/labstack/echo/v4"

func straniciExecute(c echo.Context) error {
	c.Logger().Debugf("in straniciExecute")
	lang := c.Param("lang")
	/*if lang == "" && len(c.Request().Header["Accept-Language"]) > 0 {
		lang = c.Request().Header["Accept-Language"][0]
	}
	*/
	pageAlias := c.Param("stranica")
	paragraphAlias := c.Param("celina")
	path := c.Request().URL.Path
	data := spf("Params: pagealias: %s, paragraphalias:%s,lang:%s; Path:%s; PathInfo:%s; QUERY_STRING:%s",
		pageAlias, paragraphAlias, lang, c.Path(), path, c.QueryString())
	c.Logger().Debug(data)
	return c.Render(200, "stranici/execute",
		Map{
			"title":    "Страница еди коя си!",
			"greeting": "Добре дошли! на страница " + data,
		},
	)
}
