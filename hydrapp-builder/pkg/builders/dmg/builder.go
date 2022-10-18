package dmg

import (
	"context"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/executors"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/dmg"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/rpm"
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
	appName, // Human-readable name for the app
	gpgKeyContent, // base64-encoded GPG key contents
	gpgKeyPassword string, // base64-encoded password for the GPG key
	universal bool, // Build universal (amd64 and arm64) binary instead of amd64 only
	packages []string, // MacPorts packages to install
	releases []rpm.Release, // App releases
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
		appName,
		gpgKeyContent,
		gpgKeyPassword,
		universal,
		packages,
		releases,
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
	appName,
	gpgKeyContent,
	gpgKeyPassword string
	universal bool
	packages  []string
	releases  []rpm.Release
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
			"APP_NAME":         b.appName,
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
		func(workdir string) error {
			return utils.WriteRenders(
				workdir,
				[]*renderers.Renderer{
					dmg.NewInfoRenderer(
						b.appID,
						b.appName,
						b.releases,
					),
				},
				b.overwrite,
			)
		},
	)
}
