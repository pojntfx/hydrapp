package deb

import (
	_ "embed"
	"path/filepath"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed changelog
var changelogTemplate string

type changelogData struct {
	AppID       string
	AppReleases []renderers.Release
}

func NewChangelogRenderer(
	appID string,
	appReleases []renderers.Release,
) renderers.Renderer {
	return renderers.NewRenderer(
		filepath.Join("debian", "changelog"),
		changelogTemplate,
		changelogData{appID, appReleases},
	)
}
