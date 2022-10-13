package apk

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed manifest.xml
var manifestTemplate string

type manifestData struct {
	AppID   string
	AppName string
}

func NewManifestRenderer(
	appID string,
	appName string,
) *renderers.Renderer {
	return renderers.NewRenderer("AndroidManifest.xml", manifestTemplate, manifestData{appID, appName})
}
