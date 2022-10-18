package deb

import (
	"context"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/executors"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/deb"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/rpm"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/xdg"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/utils"
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
	appID, // DEB app ID to use
	gpgKeyContent, // base64-encoded GPG key contents
	gpgKeyPassword, // base64-encoded password for the GPG key
	gpgKeyID, // ID of the GPG key to use
	baseURL, // Base URL where the repo is to be hosted
	packageVersion, // DEB package version
	os, // OS to build for
	distro, // Distro to build for
	mirrorsite string, // Mirror to use
	components []string, // Components to use
	debootstrapopts, // Options to pass to debootstrap
	architecture string, // Architecture to build for
	releases []rpm.Release, // App releases
	appDescription string, // App description
	appSummary, // App summary
	appURL, // App URL
	appGit string, // App Git repo URL
	extraDebianPackages []rpm.Package, // Extra Debian packages
	appSPDX, // App SPDX license identifier
	appLicenseDate, // App license date
	appLicenseText, // App license text
	appName string, // App name
	overwrite bool, // Overwrite files even if they exist
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
		gpgKeyContent,
		gpgKeyPassword,
		gpgKeyID,
		baseURL,
		packageVersion,
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
		appLicenseDate,
		appLicenseText,
		appName,
		overwrite,
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
	gpgKeyContent,
	gpgKeyPassword,
	gpgKeyID,
	baseURL,
	packageVersion,
	os,
	distro,
	mirrorsite string
	components []string
	debootstrapopts,
	architecture string
	releases       []rpm.Release
	appDescription string
	appSummary,
	appURL,
	appGit string
	extraDebianPackages []rpm.Package
	appSPDX,
	appLicenseDate,
	appLicenseText,
	appName string
	overwrite bool
}

func (b *Builder) Build() error {
	return executors.DockerRunImage(
		b.ctx,
		b.cli,
		b.image,
		b.pull,
		true,
		b.src,
		b.dst,
		b.onID,
		b.onOutput,
		map[string]string{
			"APP_ID":           b.appID,
			"GPG_KEY_CONTENT":  b.gpgKeyContent,
			"GPG_KEY_PASSWORD": b.gpgKeyPassword,
			"GPG_KEY_ID":       b.gpgKeyID,
			"BASE_URL":         b.baseURL,
			"OS":               b.os,
			"DISTRO":           b.distro,
			"MIRRORSITE":       b.mirrorsite,
			"COMPONENTS":       strings.Join(b.components, " "),
			"DEBOOTSTRAPOPTS":  b.debootstrapopts,
			"ARCHITECTURE":     b.architecture,
			"PACKAGE_VERSION":  b.packageVersion,
		},
		func(workdir string) error {
			return utils.WriteRenders(
				workdir,
				[]*renderers.Renderer{
					xdg.NewDesktopRenderer(
						b.appID,
						b.appName,
						b.appDescription,
					),
					xdg.NewMetainfoRenderer(
						b.appID,
						b.appName,
						b.appDescription,
						b.appSummary,
						b.appSPDX,
						b.appURL,
						b.releases,
					),
					deb.NewChangelogRenderer(
						b.appID,
						b.releases,
					),
					deb.NewCompatRenderer(),
					deb.NewFormatRenderer(),
					deb.NewOptionsRenderer(),
					deb.NewControlRenderer(
						b.appID,
						b.appDescription,
						b.appSummary,
						b.appURL,
						b.appGit,
						b.releases,
						b.extraDebianPackages,
					),
					deb.NewCopyrightRenderer(
						b.appID,
						b.appGit,
						b.appSPDX,
						b.appLicenseDate,
						b.appLicenseText,
						b.releases,
					),
					deb.NewRulesRenderer(
						b.appID,
					),
				},
				b.overwrite,
			)
		},
	)
}
