package xdg

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed metainfo.xml
var metainfoTemplate string

type metainfoData struct {
	AppID          string
	AppName        string
	AppDescription string
	AppSummary     string
	AppSPDX        string
	AppURL         string
	AppReleases    []renderers.Release
}

func NewMetainfoRenderer(
	appID string,
	appName string,
	appDescription string,
	appSummary string,
	appSPDX string,
	appURL string,
	appReleases []renderers.Release,
) *renderers.Renderer {
	return renderers.NewRenderer(
		appID+".metainfo.xml",
		metainfoTemplate,
		metainfoData{appID, appName, appDescription, appSummary, appSPDX, appURL, appReleases},
	)
}
