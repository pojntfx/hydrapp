package flatpak

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/example/pkg/executors"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-flatpak"
)

func NewBuilder(
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
) *Builder {
	return &Builder{
		ctx,
		cli,

		image,
		pull,
		dst,
		appID,
		gpgKeyContent,
		gpgKeyPassword,
		baseURL,
		architecture,
	}
}

type Builder struct {
	ctx context.Context
	cli *client.Client

	image string
	pull  bool
	dst,
	appID,
	gpgKeyContent,
	gpgKeyPassword,
	baseURL,
	architecture string
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
			"GPG_KEY_CONTENT":  b.gpgKeyContent,
			"GPG_KEY_PASSWORD": b.gpgKeyPassword,
			"BASE_URL":         b.baseURL,
			"ARCHITECTURE":     b.architecture,
		},
	)
}
