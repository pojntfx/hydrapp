package apk

//go:generate docker build -t ghcr.io/pojntfx/hydrapp-build-deb .

import (
	"context"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/example/pkg/executors"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-apk"
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
	androidCertContent, // base64-encoded Android cert contents
	androidCertPassword, // base64-encoded password for the Android cert
	baseURL string, // Base URL where the repo is to be hosted
) error {
	return executors.DockerRunImage(
		ctx,
		cli,

		image,
		pull,
		false,
		dst,
		map[string]string{
			"APP_ID":                appID,
			"GPG_KEY_CONTENT":       gpgKeyContent,
			"GPG_KEY_PASSWORD":      gpgKeyPassword,
			"ANDROID_CERT_CONTENT":  androidCertContent,
			"ANDROID_CERT_PASSWORD": androidCertPassword,
			"BASE_URL":              baseURL,
		},
	)
}
