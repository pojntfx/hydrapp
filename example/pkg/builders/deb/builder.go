package deb

//go:generate docker build -t ghcr.io/pojntfx/hydrapp-build-deb .

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/example/pkg/executors"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-deb"
)

func Build(
	ctx context.Context,
	cli *client.Client,

	image string, // OCI image to use
	pull bool, // Whether to pull the image or not
	dst, // Output directory
	appID, // DEB app ID to use
	gpgKeyContent, // base64-encoded GPG key contents
	gpgKeyPassword, // base64-encoded password for the GPG key
	gpgKeyID, // ID of the GPG key to use
	baseURL, // Base URL where the repo is to be hosted
	packageVersion, // DEB package version
	os, // OS to build for
	distro, // Distro to build for
	mirrorsite, // Mirror to use
	components, // Space-separated list of components to use
	debootstrapopts, // Options to pass to debootstrap
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
			"GPG_KEY_ID":       gpgKeyID,
			"BASE_URL":         baseURL,
			"OS":               os,
			"DISTRO":           distro,
			"MIRRORSITE":       mirrorsite,
			"COMPONENTS":       components,
			"DEBOOTSTRAPOPTS":  debootstrapopts,
			"ARCHITECTURE":     architecture,
			"PACKAGE_VERSION":  packageVersion,
		},
	)
}
