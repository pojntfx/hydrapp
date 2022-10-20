package rpm

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed spec.spec
var specTemplate string

type specData struct {
	AppID          string
	AppName        string
	AppDescription string
	AppSummary     string
	AppSPDX        string
	AppURL         string
	AppReleases    []renderers.Release
	ExtraPackages  []Package
}

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewSpecRenderer(
	appID string,
	appName string,
	appDescription string,
	appSummary string,
	appSPDX string,
	appURL string,
	appReleases []renderers.Release,
	extraPackages []Package,
) *renderers.Renderer {
	return renderers.NewRenderer(
		appID+".spec",
		specTemplate,
		specData{appID, appName, appDescription, appSummary, appSPDX, appURL, appReleases, extraPackages},
	)
}
