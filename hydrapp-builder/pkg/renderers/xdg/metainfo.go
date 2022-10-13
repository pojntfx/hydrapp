package xdg

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/rpm"
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
	AppReleases    []rpm.Release
}

func NewMetainfoRenderer(
	appID string,
	appName string,
	appDescription string,
	appSummary string,
	appSPDX string,
	appURL string,
	appReleases []rpm.Release,
) *renderers.Renderer {
	return renderers.NewRenderer(
		appID+".metainfo.xml",
		metainfoTemplate,
		metainfoData{appID, appName, appDescription, appSummary, appSPDX, appURL, appReleases},
	)
}
