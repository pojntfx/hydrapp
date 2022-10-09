package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"sync"

	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/example/pkg/builders"
	"github.com/pojntfx/hydrapp/example/pkg/builders/apk"
	"github.com/pojntfx/hydrapp/example/pkg/builders/deb"
	"github.com/pojntfx/hydrapp/example/pkg/builders/dmg"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	appID := flag.String("app-id", "com.pojtinger.felicitas.hydrapp.example", "Android app ID to use")
	appName := flag.String("app-name", "Hydrapp Example", "Human-readable name for the app")

	pull := flag.Bool("pull", false, "Whether to pull the images or not")
	dst := flag.String("dst", filepath.Join(pwd, "out"), "Output directory")
	baseURL := flag.String("base-url", "https://pojntfx.github.io/hydrapp/", "Base URL where the repos are to be hosted")

	gpgKeyContent := flag.String("gpg-key-content", "", "base64-encoded GPG key contents")
	gpgKeyPassword := flag.String("gpg-key-password", "", " base64-encoded password for the GPG key")
	gpgKeyID := flag.String("gpg-key-id", "", "ID of the GPG key to use")

	apkCertContent := flag.String("apk-cert-content", "", "base64-encoded Android cert contents")
	apkCertPassword := flag.String("apk-cert-password", "", " base64-encoded password for the Android cert")

	debPackageVersion := flag.String("deb-package-version", "0.0.1", "DEB package version")

	dmgUniversal := flag.Bool("dmg-universal", true, "Whether to build a universal instead of amd64-only binary and DMG image")

	concurrency := flag.Int("concurrency", 1, "Maximum amount of concurrent builders to run at once")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	bdrs := []builders.Builder{
		apk.NewBuilder(
			ctx,
			cli,

			apk.Image,
			*pull,
			filepath.Join(*dst, "apk"),
			*appID,
			*gpgKeyContent,
			*gpgKeyPassword,
			*apkCertContent,
			*apkCertPassword,
			*baseURL+"apk",
		),
		deb.NewBuilder(
			ctx,
			cli,

			deb.Image,
			*pull,
			filepath.Join(*dst, "deb", "debian", "sid", "x86_64"),
			*appID,
			*gpgKeyContent,
			*gpgKeyPassword,
			*gpgKeyID,
			*baseURL+"deb/debian/sid/x86_64",
			*debPackageVersion,
			"debian",
			"sid",
			"http://http.us.debian.org/debian",
			[]string{"main", "contrib"},
			"",
			"amd64",
		),
		dmg.NewBuilder(
			ctx,
			cli,

			dmg.Image,
			*pull,
			filepath.Join(*dst, "dmg"),
			*appID,
			*appName,
			*gpgKeyContent,
			*gpgKeyPassword,
			*dmgUniversal,
			[]string{},
		),
	}

	semaphore := make(chan struct{}, *concurrency)
	var wg sync.WaitGroup
	for _, b := range bdrs {
		wg.Add(1)

		semaphore <- struct{}{}

		go func(builder builders.Builder) {
			defer func() {
				<-semaphore

				wg.Done()
			}()

			if err := builder.Build(); err != nil {
				panic(err)
			}
		}(b)
	}

	wg.Wait()
}
