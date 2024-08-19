package apk

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed config.yml
var configTemplate string

type configData struct {
	AppName,
	BaseURL string
}

func NewConfigRenderer(
	appName,
	baseURL string,
) renderers.Renderer {
	return renderers.NewRenderer(
		"config.yml",
		configTemplate,
		configData{appName, baseURL},
	)
}
