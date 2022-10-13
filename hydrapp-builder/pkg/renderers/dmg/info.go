package dmg

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/rpm"
)

//go:embed info.plist
var infoTemplate string

type infoData struct {
	AppID       string
	AppName     string
	AppReleases []rpm.Release
}

func NewInfoRenderer(
	appID string,
	appName string,
	appReleases []rpm.Release,
) *renderers.Renderer {
	return renderers.NewRenderer(
		"Info.plist",
		infoTemplate,
		infoData{appID, appName, appReleases},
	)
}
