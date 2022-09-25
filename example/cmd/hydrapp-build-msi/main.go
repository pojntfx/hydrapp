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

	image := flag.String("image", "ghcr.io/pojntfx/hydrapp-build-msi", "OCI image to use")
	pull := flag.Bool("pull", true, "Whether to pull the image or not")
	dst := flag.String("dst", filepath.Join(pwd, "out", "msi"), "Output directory")
	appID := flag.String("app-id", "com.pojtinger.felicitas.hydrapp.example", "Windows app ID to use")
	gpgKeyContent := flag.String("gpg-key-content", "", "base64-encoded GPG key contents")
	gpgKeyPassword := flag.String("gpg-key-password", "", " base64-encoded password for the GPG key")
	architectures := flag.String("architectures", "386 amd64 arm arm64", "Space-separated list of architectures to build for")

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
		*dst,
		map[string]string{
			"APP_ID":           *appID,
			"GPG_KEY_CONTENT":  *gpgKeyContent,
			"GPG_KEY_PASSWORD": *gpgKeyPassword,
			"ARCHITECTURES":    *architectures,
		},
	); err != nil {
		panic(err)
	}
}
