package main

import (
	"context"
	"flag"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/example/pkg/executors"
)

func main() {
	image := flag.String("image", "ghcr.io/pojntfx/hydrapp-build-rpm", "OCI image to use")
	pull := flag.Bool("pull", true, "Whether to pull the image or not")
	dst := flag.String("dst", "out", "Output directory")
	appID := flag.String("app-id", "com.pojtinger.felicitas.hydrapp.example", "RPM app ID to use")
	gpgKeyContent := flag.String("gpg-key-content", "", "base64-encoded GPG key contents")
	gpgKeyPassword := flag.String("gpg-key-password", "", " base64-encoded password for the GPG key")
	gpgKeyID := flag.String("gpg-key-id", "", "ID of the GPG key to use")
	baseURL := flag.String("base-url", "https://pojntfx.github.io/hydrapp", "Base URL where the repo is to be hosted")
	distro := flag.String("distro", "fedora-36", "Distro to build for")
	architecture := flag.String("architecture", "amd64", "Architecture to build for")
	packageVersion := flag.String("package-version", "0.0.1", "RPM package version")
	packageSuffix := flag.String("package-suffix", "1.fc36", "RPM package suffix")

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
			"DISTRO":           *distro,
			"ARCHITECTURE":     *architecture,
			"PACKAGE_VERSION":  *packageVersion,
			"PACKAGE_SUFFIX":   *packageSuffix,
		},
	); err != nil {
		panic(err)
	}
}
