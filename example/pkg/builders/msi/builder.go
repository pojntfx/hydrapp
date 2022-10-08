package msi

//go:generate docker build -t ghcr.io/pojntfx/hydrapp-build-msi .

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/example/pkg/executors"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-msi"
)

func Build(
	ctx context.Context,
	cli *client.Client,

	image string, // OCI image to use
	pull bool, // Whether to pull the image or not
	dst, // Output directory
	appID, // Android app ID to use
	appName, // Human-readable name for the app,
	gpgKeyContent, // base64-encoded GPG key contents
	gpgKeyPassword, // base64-encoded password for the GPG key
	architecture, // Architecture to build for
	packages string, // Space-separated list of MSYS2 packages to install. Only supported for amd64.
) error {
	return executors.DockerRunImage(
		ctx,
		cli,
		image,
		pull,
		false,
		dst,
		map[string]string{
			"APP_ID":           appID,
			"APP_NAME":         appName,
			"GPG_KEY_CONTENT":  gpgKeyContent,
			"GPG_KEY_PASSWORD": gpgKeyPassword,
			"ARCHITECTURE":     architecture,
			"MSYS2PACKAGES":    packages,
		},
	)
}
