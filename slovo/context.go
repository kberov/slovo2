package slovo

import (
	"github.com/jmoiron/sqlx"
	"github.com/kberov/slovo2/model"
	"github.com/labstack/echo/v4"
)

type Context struct {
	echo.Context
	StraniciArgs *model.StraniciArgs
}

func (c *Context) DB() *sqlx.DB {
	return model.DB()
}

func (c *Context) BindArgs() (*model.StraniciArgs, error) {
	if c.StraniciArgs.UserID > 0 {
		return c.StraniciArgs, nil
	}
	err := c.Bind(c.StraniciArgs)
	return c.StraniciArgs, err
}

func slovoContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(&Context{c, new(model.StraniciArgs)})
	}
}
