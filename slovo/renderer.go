package slovo

import (
	"io"

	"github.com/kberov/gledki"
	"github.com/labstack/echo/v4"
)

type EchoRenderer struct {
	*gledki.Gledki
}

func GledkiMust(root string, ext string, tags [2]string, loadFiles bool, logger gledki.Logger) *EchoRenderer {
	tpls, err := gledki.New(root, ext, tags, false)
	if err != nil {
		logger.Fatal(err.Error())
	}
	tpls.Logger = logger
	r := &EchoRenderer{tpls}
	return r
}

// Render abides to the echo.Echo interface for echo.Renderer, but expects the
// template data to be of type gledki.Stash which actually is map[string]any.
func (g *EchoRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	stash, isStash := data.(gledki.Stash)
	if isStash {
		g.MergeStash(stash)
	} else {
		c.Echo().Logger.Warn("'data' parameter must be of type gledki.Stash for the GledkiRenderer() to replace values in teplates.")
	}
	_, err := g.Execute(w, name)
	return err
}
