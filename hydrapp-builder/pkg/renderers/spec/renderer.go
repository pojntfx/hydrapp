package spec

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed rpm.spec
var template string

type data struct {
	AppID             string
	AppName           string
	AppDescription    string
	AppSummary        string
	AppSPDX           string
	AppURL            string
	AppReleases       []Release
	ExtraRHELPackages []Package
	ExtraSUSEPackages []Package
}

type Release struct {
	Version     string `json:"version"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Email       string `json:"email"`
}

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func NewRenderer(
	appID string,
	appName string,
	appDescription string,
	appSummary string,
	appSPDX string,
	appURL string,
	appReleases []Release,
	extraRHELPackages []Package,
	extraSUSEPackages []Package,
) *renderers.Renderer {
	return renderers.NewRenderer(
		appID+".spec",
		template,
		data{appID, appName, appDescription, appSummary, appSPDX, appURL, appReleases, extraRHELPackages, extraSUSEPackages},
	)
}
