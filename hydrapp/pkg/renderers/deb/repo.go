package deb

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed repo.conf
var repoTemplate string

type repoData struct {
	AppName string
}

func NewRepoRenderer(
	appName string,
) renderers.Renderer {
	return renderers.NewRenderer(
		"repo.conf",
		repoTemplate,
		repoData{appName},
	)
}
