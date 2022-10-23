package dmg

import (
	"context"
	"encoding/base64"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/builders"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/executors"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/dmg"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/utils"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-dmg"
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
	appID, // macOS app ID to use
	appName string, // Human-readable name for the app
	gpgKeyContent []byte, // GPG key contents
	gpgKeyPassword string, // Password for the GPG key
	universal bool, // Build universal (amd64 and arm64) binary instead of amd64 only
	packages []string, // MacPorts packages to install
	releases []renderers.Release, // App releases
	overwrite bool, // Overwrite files even if they exist
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
		appName,
		base64.StdEncoding.EncodeToString(gpgKeyContent),
		base64.StdEncoding.EncodeToString([]byte(gpgKeyPassword)),
		universal,
		packages,
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
	appName,
	gpgKeyContent,
	gpgKeyPassword string
	universal bool
	packages  []string
	releases  []renderers.Release
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
			dmg.NewInfoRenderer(
				appID,
				appName,
				b.releases,
			),
		},
		b.overwrite,
	)
}

func (b *Builder) Build() error {
	dst := b.dst
	appID := b.appID
	appName := b.appName

	if b.unstable {
		dst = filepath.Join(dst, builders.UnstablePathSuffix)
		appID += builders.UnstableIDSuffix
		appName += builders.UnstableNameSuffix
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
			"APP_NAME":         appName,
			"GPG_KEY_CONTENT":  b.gpgKeyContent,
			"GPG_KEY_PASSWORD": b.gpgKeyPassword,
			"ARCHITECTURES": func() string {
				if b.universal {
					return "amd64 arm64"
				}

				return "amd64"
			}(),
			"MACPORTS": strings.Join(b.packages, " "),
		},
		b.Render,
	)
}
