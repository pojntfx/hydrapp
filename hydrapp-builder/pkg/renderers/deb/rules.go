package deb

import (
	_ "embed"
	"path/filepath"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed rules
var rulesTemplate string

type rulesData struct {
	AppID string
}

func NewRulesRenderer(
	appID string,
) *renderers.Renderer {
	return renderers.NewRenderer(
		filepath.Join("debian", "rules"),
		rulesTemplate,
		rulesData{appID},
	)
}
