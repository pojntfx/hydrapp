package deb

import (
	_ "embed"
	"path/filepath"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/rpm"
)

//go:embed control
var controlTemplate string

type controlData struct {
	AppID          string
	AppDescription string
	AppSummary     string
	AppURL         string
	AppGit         string
	AppReleases    []rpm.Release
	ExtraPackages  []rpm.Package
}

func NewControlRenderer(
	appID string,
	appDescription string,
	appSummary string,
	appURL string,
	appGit string,
	appReleases []rpm.Release,
	extraPackages []rpm.Package,
) *renderers.Renderer {
	return renderers.NewRenderer(
		filepath.Join("debian", "control"),
		controlTemplate,
		controlData{appID, appDescription, appSummary, appURL, appGit, appReleases, extraPackages},
	)
}
