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

	image := flag.String("image", "ghcr.io/pojntfx/hydrapp-build-apk", "OCI image to use")
	pull := flag.Bool("pull", true, "Whether to pull the image or not")
	dst := flag.String("dst", filepath.Join(pwd, "out", "apk"), "Output directory")
	appID := flag.String("app-id", "com.pojtinger.felicitas.hydrapp.example", "Android app ID to use")
	gpgKeyContent := flag.String("gpg-key-content", "", "base64-encoded GPG key contents")
	gpgKeyPassword := flag.String("gpg-key-password", "", " base64-encoded password for the GPG key")
	androidCertContent := flag.String("android-cert-content", "", "base64-encoded Android cert contents")
	androidCertPassword := flag.String("android-cert-password", "", " base64-encoded password for the Android cert")
	baseURL := flag.String("base-url", "https://pojntfx.github.io/hydrapp/fdroid", "Base URL where the repo is to be hosted")

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
		false,
		*dst,
		map[string]string{
			"APP_ID":                *appID,
			"GPG_KEY_CONTENT":       *gpgKeyContent,
			"GPG_KEY_PASSWORD":      *gpgKeyPassword,
			"ANDROID_CERT_CONTENT":  *androidCertContent,
			"ANDROID_CERT_PASSWORD": *androidCertPassword,
			"BASE_URL":              *baseURL,
		},
	); err != nil {
		panic(err)
	}
}
