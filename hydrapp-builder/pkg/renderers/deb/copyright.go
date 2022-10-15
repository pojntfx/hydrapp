package deb

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/rpm"
)

//go:embed copyright
var copyrightTemplate string

type copyrightData struct {
	AppID          string
	AppGit         string
	AppSPDX        string
	AppLicenseDate string
	AppLicenseText string
	AppReleases    []rpm.Release
}

func NewCopyrightRenderer(
	appID string,
	appGit string,
	appSPDX string,
	appLicenseDate string,
	appLicenseText string,
	appReleases []rpm.Release,
) *renderers.Renderer {
	return renderers.NewRenderer(
		"copyright",
		copyrightTemplate,
		copyrightData{appID, appGit, appSPDX, appLicenseDate, appLicenseText, appReleases},
	)
}
