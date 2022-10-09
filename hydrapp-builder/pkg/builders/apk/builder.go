package apk

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/executors"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-apk"
)

func NewBuilder(
	ctx context.Context,
	cli *client.Client,

	image string, // OCI image to use
	pull bool, // Whether to pull the image or not
	src, // Input directory
	dst string, // Output directory
	onID func(id string), // Callback to handle container ID
	appID, // Android app ID to use
	gpgKeyContent, // base64-encoded GPG key contents
	gpgKeyPassword, // base64-encoded password for the GPG key
	androidCertContent, // base64-encoded Android cert contents
	androidCertPassword, // base64-encoded password for the Android cert
	baseURL string, // Base URL where the repo is to be hosted
) *Builder {
	return &Builder{
		ctx,
		cli,

		image,
		pull,
		src,
		dst,
		onID,
		appID,
		gpgKeyContent,
		gpgKeyPassword,
		androidCertContent,
		androidCertPassword,
		baseURL,
	}
}

type Builder struct {
	ctx context.Context
	cli *client.Client

	image string
	pull  bool
	src,
	dst string
	onID func(id string)
	appID,
	gpgKeyContent,
	gpgKeyPassword,
	androidCertContent,
	androidCertPassword,
	baseURL string
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
		map[string]string{
			"APP_ID":                b.appID,
			"GPG_KEY_CONTENT":       b.gpgKeyContent,
			"GPG_KEY_PASSWORD":      b.gpgKeyPassword,
			"ANDROID_CERT_CONTENT":  b.androidCertContent,
			"ANDROID_CERT_PASSWORD": b.androidCertPassword,
			"BASE_URL":              b.baseURL,
		},
	)
}
