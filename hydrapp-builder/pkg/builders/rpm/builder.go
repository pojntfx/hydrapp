package rpm

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/executors"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-rpm"
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
	appID, // RPM app ID to use
	gpgKeyContent, // base64-encoded GPG key contents
	gpgKeyPassword, // base64-encoded password for the GPG key
	gpgKeyID, // ID of the GPG key to use
	baseURL, // Base URL where the repo is to be hosted
	packageVersion, // RPM package version
	distro, // Distro to build for
	architecture, // Architecture to build for
	packageSuffix string, // RPM package suffix
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
		gpgKeyContent,
		gpgKeyPassword,
		gpgKeyID,
		baseURL,
		packageVersion,
		distro,
		architecture,
		packageSuffix,
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
	packageVersion,
	distro,
	architecture,
	packageSuffix string
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
			"GPG_KEY_CONTENT":  b.gpgKeyContent,
			"GPG_KEY_PASSWORD": b.gpgKeyPassword,
			"GPG_KEY_ID":       b.gpgKeyID,
			"BASE_URL":         b.baseURL,
			"DISTRO":           b.distro,
			"ARCHITECTURE":     b.architecture,
			"PACKAGE_VERSION":  b.packageVersion,
			"PACKAGE_SUFFIX":   b.packageSuffix,
		},
		func(workdir string) error {
			return nil
		},
	)
}
