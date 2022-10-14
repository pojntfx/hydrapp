package deb

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed options
var optionsTemplate string

type optionsData struct{}

func NewOptionsRenderer() *renderers.Renderer {
	return renderers.NewRenderer("source/options", optionsTemplate, optionsData{})
}
