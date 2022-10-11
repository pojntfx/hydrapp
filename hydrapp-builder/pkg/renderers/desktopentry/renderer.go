package desktopentry

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed entry.desktop
var template string

type data struct {
	AppID          string
	AppName        string
	AppDescription string
}

func NewRenderer(
	appID string,
	appName string,
	appDescription string,
) *renderers.Renderer {
	return renderers.NewRenderer(appID+".desktop", template, data{appID, appName, appDescription})
}
