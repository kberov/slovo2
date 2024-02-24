package slovo

import (
	"io"

	"github.com/kberov/gledki"
	"github.com/labstack/echo/v4"
)

type EchoRenderer struct {
	*gledki.Gledki
}

func GledkiMust(roots []string, ext string, tags [2]string, loadFiles bool, logger gledki.Logger) *EchoRenderer {
	gledki.CacheTemplates = !Cfg.Debug
	// logger.Debugf("CacheTemplates: %v", gledki.CacheTemplates)
	tpls, err := gledki.New(roots, ext, tags, false)
	if err != nil {
		logger.Fatal(err.Error())
	}
	tpls.Logger = logger
	return &EchoRenderer{tpls}
}

// Render abides to the echo.Echo interface for echo.Renderer, but expects the
// template data to be of type gledki.Stash which actually is map[string]any.
func (g *EchoRenderer) Render(w io.Writer, name string, data any, c echo.Context) error {
	if stash, isStash := data.(gledki.Stash); isStash {
		g.MergeStash(stash)
	} else {
		c.Logger().Fatal("'data' parameter must be of type gledki.Stash for the GledkiRenderer() to interpolate values in templates.")
	}
	_, err := g.Execute(w, name)
	return err
}
