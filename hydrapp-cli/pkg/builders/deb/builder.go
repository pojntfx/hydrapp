package deb

import (
	"context"
	"encoding/base64"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/builders"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/executors"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/renderers/deb"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/renderers/rpm"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/renderers/xdg"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/utils"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-deb"
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
	appID string, // DEB app ID to use
	pgpKeyContent []byte, // PGP key contents
	pgpKeyPassword, // password for the PGP key
	pgpKeyID, // ID of the PGP key to use
	baseURL, // Base URL where the repo is to be hosted
	os, // OS to build for
	distro, // Distro to build for
	mirrorsite string, // Mirror to use
	components []string, // Components to use
	debootstrapopts, // Options to pass to debootstrap
	architecture string, // Architecture to build for
	releases []renderers.Release, // App releases
	appDescription string, // App description
	appSummary, // App summary
	appURL, // App URL
	appGit string, // App Git repo URL
	extraDebianPackages []rpm.Package, // Extra Debian packages
	appSPDX, // App SPDX license identifier
	appLicenseText, // App license text
	appName string, // App name
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
		os,
		distro,
		mirrorsite,
		components,
		debootstrapopts,
		architecture,
		releases,
		appDescription,
		appSummary,
		appURL,
		appGit,
		extraDebianPackages,
		appSPDX,
		strings.Replace(appLicenseText, "\n\n", "\n.\n", -1),
		appName,
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
	os,
	distro,
	mirrorsite string
	components []string
	debootstrapopts,
	architecture string
	releases       []renderers.Release
	appDescription string
	appSummary,
	appURL,
	appGit string
	extraDebianPackages []rpm.Package
	appSPDX,
	appLicenseText,
	appName string
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
			deb.NewChangelogRenderer(
				appID,
				b.releases,
			),
			deb.NewCompatRenderer(),
			deb.NewFormatRenderer(),
			deb.NewOptionsRenderer(
				b.goMain,
			),
			deb.NewControlRenderer(
				appID,
				b.appDescription,
				b.appSummary,
				b.appURL,
				b.appGit,
				b.releases,
				b.extraDebianPackages,
			),
			deb.NewCopyrightRenderer(
				appID,
				b.appGit,
				b.appSPDX,
				b.appLicenseText,
				b.releases,
			),
			deb.NewRulesRenderer(
				appID,
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
	baseURL := builders.GetPathForBranch(b.baseURL, b.branchID, "")

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
			"OS":               b.os,
			"DISTRO":           b.distro,
			"MIRRORSITE":       b.mirrorsite,
			"COMPONENTS":       strings.Join(b.components, " "),
			"DEBOOTSTRAPOPTS":  b.debootstrapopts,
			"ARCHITECTURE":     b.architecture,
			"PACKAGE_VERSION":  b.releases[len(b.releases)-1].Version,
			"GOMAIN":           b.goMain,
		},
		b.Render,
		[]string{},
	)
}
