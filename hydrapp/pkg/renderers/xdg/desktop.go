package xdg

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed desktop.desktop
var desktopTemplate string

type desktopData struct {
	AppID          string
	AppName        string
	AppDescription string
}

func NewDesktopRenderer(
	appID string,
	appName string,
	appDescription string,
) *renderers.Renderer {
	return renderers.NewRenderer(appID+".desktop", desktopTemplate, desktopData{appID, appName, appDescription})
}
