package slovo

import (
	"io"

	"github.com/kberov/gledki"
	"github.com/labstack/echo/v4"
)

// EchoRenderer implements the [echo.Rendere] interface by embedding
// [gledki.Gledki] and implmenting the method [EchoRenderer.Render].
type EchoRenderer struct {
	*gledki.Gledki
}

// GledkiMust instantiates [EchoRenderer]. Stops the appliaction if it happen
// to fail with a Fatal message in the log.
func GledkiMust(roots []string, ext string, tags [2]string, loadFiles bool, logger gledki.Logger) *EchoRenderer {
	_ = loadFiles // FIXME: Research and see what I had in mind with this variable.
	gledki.CacheTemplates = true
	logger.Debugf("CacheTemplates: %v", gledki.CacheTemplates)
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
	stash, isStash := data.(gledki.Stash)
	if !isStash {
		c.Logger().Fatal(
			"'data' parameter must be of type gledki.Stash for the GledkiRenderer() to interpolate values in templates.")
	}
	g.MergeStash(stash)
	_, err := g.Execute(w, name)
	return err
}
