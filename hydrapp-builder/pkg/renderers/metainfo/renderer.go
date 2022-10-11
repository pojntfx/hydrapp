package metainfo

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/spec"
)

//go:embed metainfo.xml
var template string

type data struct {
	AppID          string
	AppName        string
	AppDescription string
	AppSummary     string
	AppSPDX        string
	AppURL         string
	AppReleases    []spec.Release
}

func NewRenderer(
	appID string,
	appName string,
	appDescription string,
	appSummary string,
	appSPDX string,
	appURL string,
	appReleases []spec.Release,
) *renderers.Renderer {
	return renderers.NewRenderer(
		appID+".metainfo.xml",
		template,
		data{appID, appName, appDescription, appSummary, appSPDX, appURL, appReleases},
	)
}
