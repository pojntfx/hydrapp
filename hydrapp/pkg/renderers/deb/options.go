package deb

import (
	_ "embed"
	"path/filepath"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed options
var optionsTemplate string

type optionsData struct {
	GoMain string
}

func NewOptionsRenderer(
	goMain string,
) renderers.Renderer {
	return renderers.NewRenderer(filepath.Join("debian", "source", "options"), optionsTemplate, optionsData{goMain})
}
