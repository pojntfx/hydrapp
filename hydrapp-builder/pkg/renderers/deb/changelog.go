package deb

import (
	_ "embed"
	"path/filepath"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/rpm"
)

//go:embed changelog
var changelogTemplate string

type changelogData struct {
	AppID       string
	AppReleases []rpm.Release
}

func NewChangelogRenderer(
	appID string,
	appReleases []rpm.Release,
) *renderers.Renderer {
	return renderers.NewRenderer(
		filepath.Join("debian", "changelog"),
		changelogTemplate,
		changelogData{appID, appReleases},
	)
}
