package rpm

import (
	_ "embed"
	"time"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed spec.spec
var specTemplate string

type specData struct {
	AppID               string
	AppName             string
	AppDescription      string
	AppSummary          string
	AppSPDX             string
	AppURL              string
	AppReleases         []renderers.Release
	ExtraPackages       []Package
	GoMain              string
	GoFlags             string
	GoGenerate          string
	BranchTimestampUNIX int64
}

type Package struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
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
	goMain string,
	goFlags string,
	goGenerate string,
	branchTimestamp time.Time,
) renderers.Renderer {
	return renderers.NewRenderer(
		appID+".spec",
		specTemplate,
		specData{appID, appName, appDescription, appSummary, appSPDX, appURL, appReleases, extraPackages, goMain, goFlags, goGenerate, branchTimestamp.Unix()},
	)
}
