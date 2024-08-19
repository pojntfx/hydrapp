package flatpak

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed hydrapp.flatpakrepo
var repoTemplate string

type repoData struct {
	AppName string
	AppURL  string
	BaseURL string
}

func NewRepoRenderer(
	appName,
	appURL,
	baseURL string,
) renderers.Renderer {
	return renderers.NewRenderer(
		"hydrapp.flatpakrepo",
		repoTemplate,
		repoData{appName, appURL, baseURL},
	)
}
