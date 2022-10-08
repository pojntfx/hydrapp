package dmg

//go:generate docker build -t ghcr.io/pojntfx/hydrapp-build-dmg .

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/example/pkg/executors"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-dmg"
)

func Build(
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
	packages string, // Space-separated list of MacPorts packages to install
) error {
	return executors.DockerRunImage(
		ctx,
		cli,
		image,
		pull,
		true,
		dst,
		map[string]string{
			"APP_ID":           appID,
			"APP_NAME":         appName,
			"GPG_KEY_CONTENT":  gpgKeyContent,
			"GPG_KEY_PASSWORD": gpgKeyPassword,
			"ARCHITECTURES": func() string {
				if universal {
					return "amd64 arm64"
				}

				return "amd64"
			}(),
			"MACPORTS": packages,
		},
	)
}
