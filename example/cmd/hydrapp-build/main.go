package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/example/pkg/builders"
	"github.com/pojntfx/hydrapp/example/pkg/builders/apk"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	pull := flag.Bool("pull", false, "Whether to pull the images or not")
	dst := flag.String("dst", filepath.Join(pwd, "out"), "Output directory")
	appID := flag.String("app-id", "com.pojtinger.felicitas.hydrapp.example", "Android app ID to use")
	gpgKeyContent := flag.String("gpg-key-content", "", "base64-encoded GPG key contents")
	gpgKeyPassword := flag.String("gpg-key-password", "", " base64-encoded password for the GPG key")
	androidCertContent := flag.String("android-cert-content", "", "base64-encoded Android cert contents")
	androidCertPassword := flag.String("android-cert-password", "", " base64-encoded password for the Android cert")
	baseURL := flag.String("base-url", "https://pojntfx.github.io/hydrapp/", "Base URL where the repos are to be hosted")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	for _, builder := range []builders.Builder{
		apk.NewBuilder(
			ctx,
			cli,

			apk.Image,
			*pull,
			filepath.Join(*dst, "apk"),
			*appID,
			*gpgKeyContent,
			*gpgKeyPassword,
			*androidCertContent,
			*androidCertPassword,
			*baseURL+"apk",
		),
	} {
		if err := builder.Build(); err != nil {
			panic(err)
		}
	}
}
