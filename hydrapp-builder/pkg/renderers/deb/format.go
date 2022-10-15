package deb

import (
	_ "embed"
	"path/filepath"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed format
var formatTemplate string

type formatData struct{}

func NewFormatRenderer() *renderers.Renderer {
	return renderers.NewRenderer(filepath.Join("debian", "source", "format"), formatTemplate, formatData{})
}
