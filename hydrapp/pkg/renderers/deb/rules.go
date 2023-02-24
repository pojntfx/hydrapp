package deb

import (
	_ "embed"
	"path/filepath"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed rules
var rulesTemplate string

type rulesData struct {
	AppID      string
	GoMain     string
	GoFlags    string
	GoGenerate string
}

func NewRulesRenderer(
	appID string,
	goMain string,
	goFlags string,
	goGenerate string,
) *renderers.Renderer {
	return renderers.NewRenderer(
		filepath.Join("debian", "rules"),
		rulesTemplate,
		rulesData{appID, goMain, goFlags, goGenerate},
	)
}
