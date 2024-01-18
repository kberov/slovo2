package slovo

import "github.com/labstack/echo/v4"

func celiniExecute(c echo.Context) error {
	c.Logger().Debugf("in celiniExecute")
	pageAlias := c.Param("stranica")
	paragraphAlias := c.Param("celina")
	lang := c.Param("lang")
	path := c.Request().URL.Path
	data := spf("Params: pagealias: %s, paragraphalias:%s,lang:%s; Path:%s; PathInfo:%s; QUERY_STRING:%s",
		pageAlias, paragraphAlias, lang, c.Path(), path, c.QueryString())
	c.Logger().Debug(data)
	return c.Render(200, "hello",
		Map{
			"title":    "Целина в страница еди коя си!",
			"greeting": "Добре дошли! на страница " + data,
		},
	)
}
