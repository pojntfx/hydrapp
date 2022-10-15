package msi

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	rpm "github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/rpm"
)

//go:embed wix.wxl
var wixTemplate string

type wixData struct {
	AppID       string
	AppName     string
	AppReleases []rpm.Release
}

func NewWixRenderer(
	appID string,
	appName string,
	appReleases []rpm.Release,
) *renderers.Renderer {
	return renderers.NewRenderer(
		appID+".wxl",
		wixTemplate,
		wixData{appID, appName, appReleases},
	)
}
