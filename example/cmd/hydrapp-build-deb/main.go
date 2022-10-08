package main

import (
	"context"
	"flag"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/example/pkg/executors"
)

func main() {
	image := flag.String("image", "ghcr.io/pojntfx/hydrapp-build-deb", "OCI image to use")
	pull := flag.Bool("pull", true, "Whether to pull the image or not")
	dst := flag.String("dst", "out", "Output directory")
	appID := flag.String("app-id", "com.pojtinger.felicitas.hydrapp.example", "DEB app ID to use")
	gpgKeyContent := flag.String("gpg-key-content", "", "base64-encoded GPG key contents")
	gpgKeyPassword := flag.String("gpg-key-password", "", " base64-encoded password for the GPG key")
	gpgKeyID := flag.String("gpg-key-id", "", "ID of the GPG key to use")
	baseURL := flag.String("base-url", "https://pojntfx.github.io/hydrapp", "Base URL where the repo is to be hosted")
	packageVersion := flag.String("package-version", "0.0.1", "DEB package version")
	os := flag.String("os", "debian", "OS to build for")
	distro := flag.String("distro", "bullseye", "Distro to build for")
	mirrorsite := flag.String("mirrorsite", "http://http.us.debian.org/debian", "Mirror to use")
	components := flag.String("components", "main contrib", "Space-separated list of components to use")
	debootstrapopts := flag.String("debootstrapopts", "", "Options to pass to debootstrap")
	architecture := flag.String("architecture", "amd64", "Architecture to build for")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	if err := executors.DockerRunImage(
		ctx,
		cli,
		*image,
		*pull,
		true,
		*dst,
		map[string]string{
			"APP_ID":           *appID,
			"GPG_KEY_CONTENT":  *gpgKeyContent,
			"GPG_KEY_PASSWORD": *gpgKeyPassword,
			"GPG_KEY_ID":       *gpgKeyID,
			"BASE_URL":         *baseURL,
			"OS":               *os,
			"DISTRO":           *distro,
			"MIRRORSITE":       *mirrorsite,
			"COMPONENTS":       *components,
			"DEBOOTSTRAPOPTS":  *debootstrapopts,
			"ARCHITECTURE":     *architecture,
			"PACKAGE_VERSION":  *packageVersion,
		},
	); err != nil {
		panic(err)
	}
}
