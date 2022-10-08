package flatpak

//go:generate docker build -t ghcr.io/pojntfx/hydrapp-build-flatpak .

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/example/pkg/executors"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-flatpak"
)

func Build(
	ctx context.Context,
	cli *client.Client,

	image string, // OCI image to use
	pull bool, // Whether to pull the image or not
	dst, // Output directory
	appID, // Android app ID to use
	gpgKeyContent, // base64-encoded GPG key contents
	gpgKeyPassword, // base64-encoded password for the GPG key
	baseURL, // Base URL where the repo is to be hosted
	architecture string, // Architecture to build for
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
			"GPG_KEY_CONTENT":  gpgKeyContent,
			"GPG_KEY_PASSWORD": gpgKeyPassword,
			"BASE_URL":         baseURL,
			"ARCHITECTURE":     architecture,
		},
	)
}
