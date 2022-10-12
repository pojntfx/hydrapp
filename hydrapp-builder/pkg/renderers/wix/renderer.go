package wix

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/spec"
)

//go:embed app.wxl
var template string

type data struct {
	AppID          string
	AppName        string
	AppDescription string
	AppReleases    []spec.Release
}

func NewRenderer(
	appID string,
	appName string,
	appDescription string,
	appReleases []spec.Release,
) *renderers.Renderer {
	return renderers.NewRenderer(
		appID+".wxl",
		template,
		data{appID, appName, appDescription, appReleases},
	)
}
