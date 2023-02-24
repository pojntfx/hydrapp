package msi

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed wix.wxl
var wixTemplate string

type wixData struct {
	AppID       string
	AppName     string
	AppReleases []renderers.Release
}

func NewWixRenderer(
	appID string,
	appName string,
	appReleases []renderers.Release,
) *renderers.Renderer {
	return renderers.NewRenderer(
		appID+".wxl",
		wixTemplate,
		wixData{appID, appName, appReleases},
	)
}
