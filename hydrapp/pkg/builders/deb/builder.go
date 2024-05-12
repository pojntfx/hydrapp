package deb

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/builders"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/executors"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers/deb"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers/rpm"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers/xdg"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/utils"
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
	stdout io.Writer, // Writer to handle container output
	iconFilePath, // Path to icon to use
	appID string, // DEB app ID to use
	pgpKey []byte, // PGP key contents
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
	branchName string, // Branch Name
	branchTimestamp time.Time, // Branch timestamp
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
		stdout,
		iconFilePath,
		appID,
		base64.StdEncoding.EncodeToString(pgpKey),
		pgpKeyPassword,
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
		branchTimestamp,
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
	onID   func(id string)
	stdout io.Writer
	iconFilePath,
	appID,
	pgpKey,
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
	branchName string
	branchTimestamp time.Time
	goMain,
	goFlags,
	goGenerate string
}

func (b *Builder) Render(workdir string, ejecting bool) error {
	appID := builders.GetAppIDForBranch(b.appID, b.branchID)
	appName := builders.GetAppNameForBranch(b.appName, b.branchName)

	return utils.WriteRenders(
		filepath.Join(workdir, b.goMain),
		[]renderers.Renderer{
			xdg.NewIconRenderer(
				filepath.Join(workdir, b.goMain, b.iconFilePath),
				"icon-16x16.png",
				utils.ImageTypePNG,
				16,
				16,
			),
			xdg.NewIconRenderer(
				filepath.Join(workdir, b.goMain, b.iconFilePath),
				"icon-22x22.png",
				utils.ImageTypePNG,
				22,
				22,
			),
			xdg.NewIconRenderer(
				filepath.Join(workdir, b.goMain, b.iconFilePath),
				"icon-24x24.png",
				utils.ImageTypePNG,
				24,
				24,
			),
			xdg.NewIconRenderer(
				filepath.Join(workdir, b.goMain, b.iconFilePath),
				"icon-32x32.png",
				utils.ImageTypePNG,
				32,
				32,
			),
			xdg.NewIconRenderer(
				filepath.Join(workdir, b.goMain, b.iconFilePath),
				"icon-36x36.png",
				utils.ImageTypePNG,
				36,
				36,
			),
			xdg.NewIconRenderer(
				filepath.Join(workdir, b.goMain, b.iconFilePath),
				"icon-48x48.png",
				utils.ImageTypePNG,
				48,
				48,
			),
			xdg.NewIconRenderer(
				filepath.Join(workdir, b.goMain, b.iconFilePath),
				"icon-64x64.png",
				utils.ImageTypePNG,
				64,
				64,
			),
			xdg.NewIconRenderer(
				filepath.Join(workdir, b.goMain, b.iconFilePath),
				"icon-72x72.png",
				utils.ImageTypePNG,
				72,
				72,
			),
			xdg.NewIconRenderer(
				filepath.Join(workdir, b.goMain, b.iconFilePath),
				"icon-96x96.png",
				utils.ImageTypePNG,
				96,
				96,
			),
			xdg.NewIconRenderer(
				filepath.Join(workdir, b.goMain, b.iconFilePath),
				"icon-128x128.png",
				utils.ImageTypePNG,
				128,
				128,
			),
			xdg.NewIconRenderer(
				filepath.Join(workdir, b.goMain, b.iconFilePath),
				"icon-192x192.png",
				utils.ImageTypePNG,
				192,
				192,
			),
			xdg.NewIconRenderer(
				filepath.Join(workdir, b.goMain, b.iconFilePath),
				"icon-256x256.png",
				utils.ImageTypePNG,
				256,
				256,
			),
			xdg.NewIconRenderer(
				filepath.Join(workdir, b.goMain, b.iconFilePath),
				"icon-512x512.png",
				utils.ImageTypePNG,
				512,
				512,
			),
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
				b.branchTimestamp,
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
		b.stdout,
		map[string]string{
			"APP_ID":                appID,
			"PGP_KEY":               b.pgpKey,
			"PGP_KEY_PASSWORD":      b.pgpKeyPassword,
			"PGP_KEY_ID":            b.pgpKeyID,
			"BASE_URL":              baseURL,
			"OS":                    b.os,
			"DISTRO":                b.distro,
			"MIRRORSITE":            b.mirrorsite,
			"COMPONENTS":            strings.Join(b.components, " "),
			"DEBOOTSTRAPOPTS":       b.debootstrapopts,
			"ARCHITECTURE":          b.architecture,
			"PACKAGE_VERSION":       b.releases[0].Version,
			"GOMAIN":                b.goMain,
			"BRANCH_TIMESTAMP_UNIX": fmt.Sprintf("%v", b.branchTimestamp.Unix()),
		},
		b.Render,
		[]string{},
	)
}
