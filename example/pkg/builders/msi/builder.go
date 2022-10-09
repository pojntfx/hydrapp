package msi

import (
	"context"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/example/pkg/executors"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-msi"
)

func NewBuilder(
	ctx context.Context,
	cli *client.Client,

	image string, // OCI image to use
	pull bool, // Whether to pull the image or not
	dst, // Output directory
	appID, // Android app ID to use
	appName, // Human-readable name for the app,
	gpgKeyContent, // base64-encoded GPG key contents
	gpgKeyPassword, // base64-encoded password for the GPG key
	architecture string, // Architecture to build for
	packages []string, // MSYS2 packages to install. Only supported for amd64.
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
		architecture,
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
	gpgKeyPassword,
	architecture string
	packages []string
}

func (b *Builder) Build() error {
	return executors.DockerRunImage(
		b.ctx,
		b.cli,
		b.image,
		b.pull,
		false,
		b.dst,
		map[string]string{
			"APP_ID":           b.appID,
			"APP_NAME":         b.appName,
			"GPG_KEY_CONTENT":  b.gpgKeyContent,
			"GPG_KEY_PASSWORD": b.gpgKeyPassword,
			"ARCHITECTURE":     b.architecture,
			"MSYS2PACKAGES":    strings.Join(b.packages, " "),
		},
	)
}
