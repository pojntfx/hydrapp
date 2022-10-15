package deb

import (
	_ "embed"
	"path/filepath"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed options
var optionsTemplate string

type optionsData struct{}

func NewOptionsRenderer() *renderers.Renderer {
	return renderers.NewRenderer(filepath.Join("debian", "source", "options"), optionsTemplate, optionsData{})
}
