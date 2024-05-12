package deb

import (
	_ "embed"
	"path/filepath"
	"time"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed changelog
var changelogTemplate string

type changelogData struct {
	AppID               string
	AppReleases         []renderers.Release
	BranchTimestampUNIX int64
}

func NewChangelogRenderer(
	appID string,
	appReleases []renderers.Release,
	branchTimestamp time.Time,
) renderers.Renderer {
	return renderers.NewRenderer(
		filepath.Join("debian", "changelog"),
		changelogTemplate,
		changelogData{appID, appReleases, branchTimestamp.Unix()},
	)
}
