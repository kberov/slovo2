package slovo

import (
	"net/http"

	"github.com/kberov/slovo2/model"
	"github.com/labstack/echo/v4"
)

func celiniExecute(c echo.Context) error {
	c.Logger().Debugf("in celiniExecute")
	log := c.Logger()
	args := new(model.StraniciArgs)
	if err := c.Bind(args); err != nil {
		return c.String(http.StatusBadRequest, "Грешна заявка!")
	}
	cel := new(model.Celini)
	if err := cel.FindForDisplay(args); err != nil {
		log.Errorf("celina: %#v; error:%w; ErrType: %T; args: %#v", cel, err, err, args)
		return handleNotFound(c, args, err)
	}
	return c.Render(http.StatusOK, "hello",
		Map{
			"title":    "Целина в страница еди коя си!",
			"greeting": "Добре дошли! на страница ",
		},
	)
}
