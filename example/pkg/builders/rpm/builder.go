package rpm

//go:generate docker build -t ghcr.io/pojntfx/hydrapp-build-rpm .

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/example/pkg/executors"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-rpm"
)

func Build(
	ctx context.Context,
	cli *client.Client,

	image string, // OCI image to use
	pull bool, // Whether to pull the image or not
	dst, // Output directory
	appID, // RPM app ID to use
	gpgKeyContent, // base64-encoded GPG key contents
	gpgKeyPassword, // base64-encoded password for the GPG key
	gpgKeyID, // ID of the GPG key to use
	baseURL, // Base URL where the repo is to be hosted
	packageVersion, // RPM package version
	distro, // Distro to build for
	architecture, // Architecture to build for
	packageSuffix string, // RPM package suffix
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
			"GPG_KEY_ID":       gpgKeyID,
			"BASE_URL":         baseURL,
			"DISTRO":           distro,
			"ARCHITECTURE":     architecture,
			"PACKAGE_VERSION":  packageVersion,
			"PACKAGE_SUFFIX":   packageSuffix,
		},
	)
}
