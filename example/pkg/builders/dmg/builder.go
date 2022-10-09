package dmg

import (
	"context"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/example/pkg/executors"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-dmg"
)

func NewBuilder(
	ctx context.Context,
	cli *client.Client,

	image string, // OCI image to use
	pull bool, // Whether to pull the image or not
	dst, // Output directory
	appID, // macOS app ID to use
	appName, // Human-readable name for the app
	gpgKeyContent, // base64-encoded GPG key contents
	gpgKeyPassword string, // base64-encoded password for the GPG key
	universal bool, // Build universal (amd64 and arm64) binary instead of amd64 only
	packages []string, // MacPorts packages to install
) *Builder {
	return &Builder{
		ctx,
		cli,

		image,
		pull,
		dst,
		appID,
		appName,
		gpgKeyContent,
		gpgKeyPassword,
		universal,
		packages,
	}
}

type Builder struct {
	ctx context.Context
	cli *client.Client

	image string
	pull  bool
	dst,
	appID,
	appName,
	gpgKeyContent,
	gpgKeyPassword string
	universal bool
	packages  []string
}

func (b *Builder) Build() error {
	return executors.DockerRunImage(
		b.ctx,
		b.cli,
		b.image,
		b.pull,
		true,
		b.dst,
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
	)
}
