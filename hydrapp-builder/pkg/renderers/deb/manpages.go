package deb

import (
	_ "embed"

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
		appID+".manpages",
		manpagesTemplate,
		manpagesData{appID},
	)
}
