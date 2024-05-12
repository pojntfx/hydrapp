package apk

import (
	_ "embed"
	"time"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed manifest.xml
var manifestTemplate string

type manifestData struct {
	AppID               string
	AppName             string
	AppReleases         []renderers.Release
	BranchTimestampUNIX int64
}

func NewManifestRenderer(
	appID string,
	appName string,
	appReleases []renderers.Release,
	branchTimestamp time.Time,
) renderers.Renderer {
	return renderers.NewRenderer("AndroidManifest.xml", manifestTemplate, manifestData{appID, appName, appReleases, branchTimestamp.Unix()})
}
