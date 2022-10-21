package apk

import (
	"context"
	"encoding/base64"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/executors"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/apk"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/utils"
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
	onOutput func(shortID string, color string, timestamp int64, message string), // Callback to handle container output
	appID string, // Android app ID to use
	gpgKeyContent []byte, // GPG key contents
	gpgKeyPassword string, // Password for the GPG key
	androidCertContent []byte, // Android cert contents
	androidCertPassword string, // Password for the Android cert
	baseURL, // Base URL where the repo is to be hosted
	appName string, // App name
	overwrite bool, // Overwrite files even if they exist
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
		base64.StdEncoding.EncodeToString(gpgKeyContent),
		base64.StdEncoding.EncodeToString([]byte(gpgKeyPassword)),
		base64.StdEncoding.EncodeToString(androidCertContent),
		base64.StdEncoding.EncodeToString([]byte(androidCertPassword)),
		baseURL,
		appName,
		overwrite,
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
	androidCertContent,
	androidCertPassword,
	baseURL,
	appName string
	overwrite bool
}

func (b *Builder) Render(workdir string) error {
	return utils.WriteRenders(
		workdir,
		[]*renderers.Renderer{
			apk.NewManifestRenderer(
				b.appID,
				b.appName,
			),
			apk.NewActivityRenderer(
				b.appID,
			),
			apk.NewHeaderRenderer(),
			apk.NewImplementationRenderer(),
		},
		b.overwrite,
	)
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
			"APP_ID":                b.appID,
			"GPG_KEY_CONTENT":       b.gpgKeyContent,
			"GPG_KEY_PASSWORD":      b.gpgKeyPassword,
			"ANDROID_CERT_CONTENT":  b.androidCertContent,
			"ANDROID_CERT_PASSWORD": b.androidCertPassword,
			"BASE_URL":              b.baseURL,
		},
		b.Render,
	)
}
