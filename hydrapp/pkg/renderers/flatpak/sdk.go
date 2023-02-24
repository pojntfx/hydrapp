package flatpak

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed sdk.yaml
var sdkTemplate string

type sdkData struct{}

func NewSdkRenderer() *renderers.Renderer {
	return renderers.NewRenderer(
		"org.freedesktop.Sdk.Extension.ImageMagick.yaml",
		sdkTemplate,
		sdkData{},
	)
}
