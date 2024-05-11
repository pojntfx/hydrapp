package flatpak

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed manifest.yaml
var manifestTemplate string

type manifestData struct {
	AppID      string
	GoMain     string
	GoFlags    string
	GoGenerate string
}

func NewManifestRenderer(
	appID string,
	goMain string,
	goFlags string,
	goGenerate string,
) renderers.Renderer {
	return renderers.NewRenderer(
		appID+".yaml",
		manifestTemplate,
		manifestData{appID, goMain, goFlags, goGenerate},
	)
}
