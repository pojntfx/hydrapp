package deb

import (
	"context"
	"encoding/base64"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/builders"
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
	appID string, // DEB app ID to use
	gpgKeyContent []byte, // GPG key contents
	gpgKeyPassword, // password for the GPG key
	gpgKeyID, // ID of the GPG key to use
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
	overwrite, // Overwrite files even if they exist
	unstable bool, // Create unstable build
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
		base64.StdEncoding.EncodeToString(gpgKeyContent),
		base64.StdEncoding.EncodeToString([]byte(gpgKeyPassword)),
		gpgKeyID,
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
		appLicenseText,
		appName,
		overwrite,
		unstable,
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
	overwrite,
	unstable bool
}

func (b *Builder) Render(workdir string, ejecting bool) error {
	appID := b.appID
	appName := b.appName

	if b.unstable {
		appID += builders.UnstableIDSuffix
		appName += builders.UnstableNameSuffix
	}

	return utils.WriteRenders(
		workdir,
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
			deb.NewOptionsRenderer(),
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
			),
		},
		b.overwrite,
	)
}

func (b *Builder) Build() error {
	dst := b.dst
	appID := b.appID
	baseURL := b.baseURL

	if b.unstable {
		dst = filepath.Join(dst, builders.UnstablePathSuffix)
		appID += builders.UnstableIDSuffix
		baseURL += "/" + builders.UnstablePathSuffix
	} else {
		dst = filepath.Join(dst, builders.StablePathSuffix)
		baseURL += "/" + builders.StablePathSuffix
	}

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
			"GPG_KEY_CONTENT":  b.gpgKeyContent,
			"GPG_KEY_PASSWORD": b.gpgKeyPassword,
			"GPG_KEY_ID":       b.gpgKeyID,
			"BASE_URL":         baseURL,
			"OS":               b.os,
			"DISTRO":           b.distro,
			"MIRRORSITE":       b.mirrorsite,
			"COMPONENTS":       strings.Join(b.components, " "),
			"DEBOOTSTRAPOPTS":  b.debootstrapopts,
			"ARCHITECTURE":     b.architecture,
			"PACKAGE_VERSION":  b.releases[len(b.releases)-1].Version,
		},
		b.Render,
	)
}
