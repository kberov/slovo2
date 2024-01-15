package slovo

import "github.com/labstack/echo/v4"

func straniciExecute(c echo.Context) error {
	lang := c.Param("lang")
	if lang == "" && len(c.Request().Header["Accept-Language"]) > 0 {
		lang = c.Request().Header["Accept-Language"][0]
	}
	stash := Map{
		"title": c.Param("page_alias"),
		"lang":  lang,
	}
	return c.Render(200, "stranici/execute", stash)
}
