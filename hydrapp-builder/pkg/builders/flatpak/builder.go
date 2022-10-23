package flatpak

import (
	"context"
	"encoding/base64"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/builders"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/executors"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/flatpak"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/xdg"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/utils"
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
	gpgKeyContent []byte, // GPG key contents
	gpgKeyPassword, // Password for the GPG key
	gpgKeyID, // ID of the GPG key to use
	baseURL, // Base URL where the repo is to be hosted
	architecture, // Architecture to build for
	appName, // App name
	appDescription, // App description
	appSummary, // App summary
	appSPDX, // App SPDX license identifier
	appURL string, // App URL
	releases []renderers.Release, // App releases
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
		architecture,
		appName,
		appDescription,
		appSummary,
		appSPDX,
		appURL,
		releases,
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
	architecture,
	appName,
	appDescription,
	appSummary,
	appSPDX,
	appURL string
	releases []renderers.Release
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
			flatpak.NewManifestRenderer(
				appID,
			),
			flatpak.NewSdkRenderer(),
		},
		b.overwrite,
		ejecting,
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
			"ARCHITECTURE":     b.architecture,
		},
		b.Render,
	)
}
