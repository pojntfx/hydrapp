package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/example/pkg/executors"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	image := flag.String("image", "ghcr.io/pojntfx/hydrapp-build-flatpak", "OCI image to use")
	pull := flag.Bool("pull", true, "Whether to pull the image or not")
	dst := flag.String("dst", filepath.Join(pwd, "out", "apk"), "Output directory")
	appID := flag.String("app-id", "com.pojtinger.felicitas.hydrapp.example", "Flatpak app ID to use")
	gpgKeyContent := flag.String("gpg-key-content", "", "base64-encoded GPG key contents")
	gpgKeyPassword := flag.String("gpg-key-password", "", " base64-encoded password for the GPG key")
	baseURL := flag.String("base-url", "https://pojntfx.github.io/hydrapp/flatpak", "Base URL where the repo is to be hosted")
	architectures := flag.String("architectures", "amd64", "Space-separated list of architectures to build for")

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
			"BASE_URL":         *baseURL,
			"ARCHITECTURES":    *architectures,
		},
	); err != nil {
		panic(err)
	}
}
