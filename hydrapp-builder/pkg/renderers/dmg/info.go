package dmg

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed info.plist
var infoTemplate string

type infoData struct {
	AppID       string
	AppName     string
	AppReleases []renderers.Release
}

func NewInfoRenderer(
	appID string,
	appName string,
	appReleases []renderers.Release,
) *renderers.Renderer {
	return renderers.NewRenderer(
		"Info.plist",
		infoTemplate,
		infoData{appID, appName, appReleases},
	)
}
