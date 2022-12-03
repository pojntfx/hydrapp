package deb

import (
	_ "embed"
	"path/filepath"

	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/renderers"
)

//go:embed compat
var compatTemplate string

type compatData struct{}

func NewCompatRenderer() *renderers.Renderer {
	return renderers.NewRenderer(filepath.Join("debian", "compat"), compatTemplate, compatData{})
}
