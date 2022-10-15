package rpm

import (
	_ "embed"
	"time"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
)

//go:embed spec.spec
var specTemplate string

type specData struct {
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
	Version     string    `json:"version"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Email       string    `json:"email"`
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
	appReleases []Release,
	extraRHELPackages []Package,
	extraSUSEPackages []Package,
) *renderers.Renderer {
	return renderers.NewRenderer(
		appID+".spec",
		specTemplate,
		specData{appID, appName, appDescription, appSummary, appSPDX, appURL, appReleases, extraRHELPackages, extraSUSEPackages},
	)
}
