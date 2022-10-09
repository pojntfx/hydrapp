package deb

import (
	"context"
	"strings"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/example/pkg/executors"
)

const (
	Image = "ghcr.io/pojntfx/hydrapp-build-deb"
)

func NewBuilder(
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
	mirrorsite string, // Mirror to use
	components []string, // Components to use
	debootstrapopts, // Options to pass to debootstrap
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
		gpgKeyID,
		baseURL,
		packageVersion,
		os,
		distro,
		mirrorsite,
		components,
		debootstrapopts,
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
	gpgKeyID,
	baseURL,
	packageVersion,
	os,
	distro,
	mirrorsite string
	components []string
	debootstrapopts,
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
			"GPG_KEY_ID":       b.gpgKeyID,
			"BASE_URL":         b.baseURL,
			"OS":               b.os,
			"DISTRO":           b.distro,
			"MIRRORSITE":       b.mirrorsite,
			"COMPONENTS":       strings.Join(b.components, " "),
			"DEBOOTSTRAPOPTS":  b.debootstrapopts,
			"ARCHITECTURE":     b.architecture,
			"PACKAGE_VERSION":  b.packageVersion,
		},
	)
}
