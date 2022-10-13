package flatpak

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed manifest.yaml
var manifestTemplate string

type manifestData struct {
	AppID string
}

func NewManifestRenderer(
	appID string,
) *renderers.Renderer {
	return renderers.NewRenderer(
		appID+".yaml",
		manifestTemplate,
		manifestData{appID},
	)
}
