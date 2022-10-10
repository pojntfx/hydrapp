package androidmanifest

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed AndroidManifest.xml
var template string

type data struct {
	AppID   string
	AppName string
}

func NewRenderer(
	appID string,
	appName string,
) *renderers.Renderer[data] {
	return renderers.NewRenderer("AndroidManifest.xml", template, data{appID, appName})
}
