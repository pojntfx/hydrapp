package flatpak

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed manifest.yaml
var template string

type data struct {
	AppID string
}

func NewRenderer(
	appID string,
) *renderers.Renderer {
	return renderers.NewRenderer(
		appID+".yaml",
		template,
		data{appID},
	)
}
