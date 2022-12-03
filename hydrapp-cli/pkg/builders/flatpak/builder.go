package flatpak

import (
	"context"
	"encoding/base64"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/builders"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/executors"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/renderers/flatpak"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/renderers/xdg"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/utils"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-flatpak"
)

func NewBuilder(
	ctx context.Context,
	cli *client.Client,

	image string, // OCI image to use
	pull bool, // Whether to pull the image or not
	src, // Input directory
	dst string, // Output directory
	onID func(id string), // Callback to handle container ID
	onOutput func(shortID string, color string, timestamp int64, message string), // Callback to handle container output
	appID string, // Android app ID to use
	pgpKeyContent []byte, // PGP key contents
	pgpKeyPassword, // Password for the PGP key
	pgpKeyID, // ID of the PGP key to use
	baseURL, // Base URL where the repo is to be hosted
	architecture, // Architecture to build for
	appName, // App name
	appDescription, // App description
	appSummary, // App summary
	appSPDX, // App SPDX license identifier
	appURL string, // App URL
	releases []renderers.Release, // App releases
	overwrite bool, // Overwrite files even if they exist
	branchID, // Branch ID
	branchName, // Branch Name
	goMain, // Directory with the main package to build
	goFlags, // Flags to pass to the Go command
	goGenerate string, // Command to execute go generate with
) *Builder {
	return &Builder{
		ctx,
		cli,

		image,
		pull,
		src,
		dst,
		onID,
		onOutput,
		appID,
		base64.StdEncoding.EncodeToString(pgpKeyContent),
		base64.StdEncoding.EncodeToString([]byte(pgpKeyPassword)),
		pgpKeyID,
		baseURL,
		architecture,
		appName,
		appDescription,
		appSummary,
		appSPDX,
		appURL,
		releases,
		overwrite,
		branchID,
		branchName,
		goMain,
		goFlags,
		goGenerate,
	}
}

type Builder struct {
	ctx context.Context
	cli *client.Client

	image string
	pull  bool
	src,
	dst string
	onID     func(id string)
	onOutput func(shortID string, color string, timestamp int64, message string)
	appID,
	pgpKeyContent,
	pgpKeyPassword,
	pgpKeyID,
	baseURL,
	architecture,
	appName,
	appDescription,
	appSummary,
	appSPDX,
	appURL string
	releases  []renderers.Release
	overwrite bool
	branchID,
	branchName,
	goMain,
	goFlags,
	goGenerate string
}

func (b *Builder) Render(workdir string, ejecting bool) error {
	appID := builders.GetAppIDForBranch(b.appID, b.branchID)
	appName := builders.GetAppNameForBranch(b.appName, b.branchName)

	return utils.WriteRenders(
		filepath.Join(workdir, b.goMain),
		[]*renderers.Renderer{
			xdg.NewDesktopRenderer(
				appID,
				appName,
				b.appDescription,
			),
			xdg.NewMetainfoRenderer(
				appID,
				appName,
				b.appDescription,
				b.appSummary,
				b.appSPDX,
				b.appURL,
				b.releases,
			),
			flatpak.NewManifestRenderer(
				appID,
				b.goMain,
				b.goFlags,
				b.goGenerate,
			),
			flatpak.NewSdkRenderer(),
		},
		b.overwrite,
		ejecting,
	)
}

func (b *Builder) Build() error {
	dst := builders.GetFilepathForBranch(b.dst, b.branchID)
	appID := builders.GetAppIDForBranch(b.appID, b.branchID)
	baseURL := builders.GetPathForBranch(b.baseURL, b.branchID)

	return executors.DockerRunImage(
		b.ctx,
		b.cli,
		b.image,
		b.pull,
		true,
		b.src,
		dst,
		b.onID,
		b.onOutput,
		map[string]string{
			"APP_ID":           appID,
			"PGP_KEY_CONTENT":  b.pgpKeyContent,
			"PGP_KEY_PASSWORD": b.pgpKeyPassword,
			"PGP_KEY_ID":       b.pgpKeyID,
			"BASE_URL":         baseURL,
			"ARCHITECTURE":     b.architecture,
			"GOMAIN":           b.goMain,
		},
		b.Render,
		[]string{},
	)
}
