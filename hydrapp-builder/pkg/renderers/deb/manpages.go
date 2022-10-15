package deb

import (
	_ "embed"
	"path/filepath"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed manpages
var manpagesTemplate string

type manpagesData struct {
	AppID string
}

func NewManpagesRenderer(
	appID string,
) *renderers.Renderer {
	return renderers.NewRenderer(
		filepath.Join("debian", appID+".manpages"),
		manpagesTemplate,
		manpagesData{appID},
	)
}
