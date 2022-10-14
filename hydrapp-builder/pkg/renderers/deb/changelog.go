package deb

import (
	_ "embed"

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
		"changelog",
		changelogTemplate,
		changelogData{appID, appReleases},
	)
}
