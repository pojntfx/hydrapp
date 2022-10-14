package deb

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed format
var formatTemplate string

type formatData struct{}

func NewFormatRenderer() *renderers.Renderer {
	return renderers.NewRenderer("source/format", formatTemplate, formatData{})
}
