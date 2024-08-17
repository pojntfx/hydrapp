package apk

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed config.yml.tpl
var configTemplate string

type configData struct {
	AppName string
}

func NewConfigRenderer(
	appName string,
) renderers.Renderer {
	return renderers.NewRenderer(
		"config.yml.tpl",
		configTemplate,
		configData{appName},
	)
}
