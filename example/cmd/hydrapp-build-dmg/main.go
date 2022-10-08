package main

import (
	"context"
	"flag"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/example/pkg/executors"
)

func main() {
	image := flag.String("image", "ghcr.io/pojntfx/hydrapp-build-dmg", "OCI image to use")
	pull := flag.Bool("pull", true, "Whether to pull the image or not")
	dst := flag.String("dst", "out", "Output directory")
	appID := flag.String("app-id", "com.pojtinger.felicitas.hydrapp.example", "macOS app ID to use")
	appName := flag.String("app-name", "Hydrapp Example", "Human-readable name for the app")
	gpgKeyContent := flag.String("gpg-key-content", "", "base64-encoded GPG key contents")
	gpgKeyPassword := flag.String("gpg-key-password", "", " base64-encoded password for the GPG key")
	universal := flag.Bool("universal", true, "Build universal (amd64 and arm64) binary instead of amd64 only")
	packages := flag.String("packages", "", "Space-separated list of MacPorts packages to install")

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
			"APP_NAME":         *appName,
			"GPG_KEY_CONTENT":  *gpgKeyContent,
			"GPG_KEY_PASSWORD": *gpgKeyPassword,
			"ARCHITECTURES": func() string {
				if *universal {
					return "amd64 arm64"
				}

				return "amd64"
			}(),
			"MACPORTS": *packages,
		},
	); err != nil {
		panic(err)
	}
}
