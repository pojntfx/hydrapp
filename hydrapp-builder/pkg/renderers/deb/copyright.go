package deb

import (
	_ "embed"
	"path/filepath"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed copyright
var copyrightTemplate string

type copyrightData struct {
	AppID          string
	AppGit         string
	AppSPDX        string
	AppLicenseText string
	AppReleases    []renderers.Release
}

func NewCopyrightRenderer(
	appID string,
	appGit string,
	appSPDX string,
	appLicenseText string,
	appReleases []renderers.Release,
) *renderers.Renderer {
	return renderers.NewRenderer(
		filepath.Join("debian", "copyright"),
		copyrightTemplate,
		copyrightData{appID, appGit, appSPDX, appLicenseText, appReleases},
	)
}
