package docs

import (
	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
)

//go:embed INSTALLATION.md.tpl
var installationTemplate string

type installationData struct {
	AppID           string
	AppName         string
	AndroidRepoURL  string
	MacOSBinaryURL  string
	MacOSBinaryName string
	BinariesURL     string
	Flatpaks        []Artifact
	MSIs            []Artifact
	RPMs            []DistroArtifact
	DEBs            []DistroArtifact
	RenderDMG       bool
	RenderAPK       bool
	RenderBinaries  bool
}

type Artifact struct {
	Architecture string
	URL          string
}

type DistroArtifact struct {
	Artifact
	DistroName    string
	DistroVersion string
}

func NewInstallationRenderer(
	appID,
	appName,
	androidRepoURL,
	macOSBinaryURL,
	macOSBinaryName,
	binariesURL string,
	flatpaks,
	msis []Artifact,
	rpms,
	debs []DistroArtifact,
	renderDMG,
	renderAPK,
	renderBinaries bool,
) *renderers.Renderer {
	return renderers.NewRenderer(
		"INSTALLATION.md",
		installationTemplate,
		installationData{
			appID,
			appName,
			androidRepoURL,
			macOSBinaryURL,
			macOSBinaryName,
			binariesURL,
			flatpaks,
			msis,
			rpms,
			debs,
			renderDMG,
			renderAPK,
			renderBinaries,
		},
	)
}
