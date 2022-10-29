package rpm

import (
	"context"
	"encoding/base64"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/builders"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/executors"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/rpm"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/xdg"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/utils"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-rpm"
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
	appID string, // RPM app ID to use
	pgpKeyContent []byte, // PGP key contents
	pgpKeyPassword, // Password for the PGP key
	pgpKeyID, // ID of the PGP key to use
	baseURL, // Base URL where the repo is to be hosted
	distro, // Distro to build for
	architecture, // Architecture to build for
	packageSuffix, // RPM package suffix
	appName, // App name
	appDescription, // App description
	appSummary, // App summary
	appURL, // App URL
	appSPDX string, // App SPDX license identifier
	releases []renderers.Release, // App releases
	extraPackages []rpm.Package, // Extra RPM packages
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
		distro,
		architecture,
		packageSuffix,
		appName,
		appDescription,
		appSummary,
		appURL,
		appSPDX,
		releases,
		extraPackages,
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
	distro,
	architecture,
	packageSuffix,
	appName,
	appDescription,
	appSummary,
	appURL,
	appSPDX string
	releases      []renderers.Release
	extraPackages []rpm.Package
	overwrite     bool
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
			rpm.NewSpecRenderer(
				appID,
				appName,
				b.appDescription,
				b.appSummary,
				b.appSPDX,
				b.appURL,
				b.releases,
				b.extraPackages,
				b.goMain,
				b.goFlags,
				b.goGenerate,
			),
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
			"DISTRO":           b.distro,
			"ARCHITECTURE":     b.architecture,
			"PACKAGE_VERSION":  b.releases[len(b.releases)-1].Version,
			"PACKAGE_SUFFIX":   b.packageSuffix,
			"GOMAIN":           b.goMain,
		},
		b.Render,
	)
}
