package msi

import (
	"context"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/executors"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-msi"
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
		src,
		dst,
		onID,
		onOutput,
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
	src,
	dst string
	onID     func(id string)
	onOutput func(shortID string, color string, timestamp int64, message string)
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
		b.src,
		b.dst,
		b.onID,
		b.onOutput,
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
