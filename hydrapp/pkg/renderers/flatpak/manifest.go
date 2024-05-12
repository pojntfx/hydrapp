package flatpak

import (
	_ "embed"
	"path"
	"strings"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed manifest.yaml
var manifestTemplate string

type manifestData struct {
	AppID      string
	GoMain     string
	GoFlags    string
	GoGenerate string
	SrcDir     string
}

func NewManifestRenderer(
	appID string,
	goMain string,
	goFlags string,
	goGenerate string,
) renderers.Renderer {
	srcDir := "."
	if goMain != "." {
		goMainComponents := strings.Split(goMain, "/") // We use the UNIX file separator here since Go uses UNIX-style paths for module names

		for i := range goMainComponents {
			if i > 0 { // `goMain` always starts with "./", so skip the first folder
				srcDir = path.Join(srcDir, "..")
			}
		}
	}

	return renderers.NewRenderer(
		appID+".yaml",
		manifestTemplate,
		manifestData{appID, goMain, goFlags, goGenerate, srcDir},
	)
}
