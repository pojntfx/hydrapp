package rpm

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed hydrapp.repo
var repoTemplate string

type repoData struct {
	AppName string
	BaseURL string
}

func NewRepoRenderer(
	appName,
	baseURL string,
) renderers.Renderer {
	return renderers.NewRenderer(
		"hydrapp.repo",
		repoTemplate,
		repoData{appName, baseURL},
	)
}
