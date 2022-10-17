package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"
	"unicode"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/builders"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/builders/apk"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/builders/deb"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/builders/dmg"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/builders/flatpak"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/builders/msi"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/builders/rpm"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/utils"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	appID := flag.String("app-id", "com.pojtinger.felicitas.hydrapp.example", "Android app ID to use")
	appName := flag.String("app-name", "Hydrapp Example", "Human-readable name for the app")

	pull := flag.Bool("pull", false, "Whether to pull the images or not")
	src := flag.String("src", pwd, "Source directory")
	dst := flag.String("dst", filepath.Join(pwd, "out"), "Output directory")
	baseURL := flag.String("base-url", "https://pojntfx.github.io/hydrapp/", "Base URL where the repos are to be hosted")

	gpgKeyContent := flag.String("gpg-key-content", "", "base64-encoded GPG key contents")
	gpgKeyPassword := flag.String("gpg-key-password", "", " base64-encoded password for the GPG key")
	gpgKeyID := flag.String("gpg-key-id", "", "ID of the GPG key to use")

	apkCertContent := flag.String("apk-cert-content", "", "base64-encoded Android cert contents")
	apkCertPassword := flag.String("apk-cert-password", "", " base64-encoded password for the Android cert")

	debPackageVersion := flag.String("deb-package-version", "0.0.1", "DEB package version")

	dmgUniversal := flag.Bool("dmg-universal", true, "Whether to build a universal instead of amd64-only binary and DMG image")

	rpmPackageVersion := flag.String("rpm-package-version", "0.0.1", "RPM package version")
	rpmPackageSuffix := flag.String("rpm-package-suffix", "1.fc36", "RPM package suffix")

	concurrency := flag.Int("concurrency", 1, "Maximum amount of concurrent builders to run at once")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	handleID := func(id string) {
		s := make(chan os.Signal)
		signal.Notify(s, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-s

			log.Println("Gracefully shutting down")

			go func() {
				<-s

				log.Println("Forcing shutdown")

				os.Exit(1)
			}()

			if err := cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{
				Force: true,
			}); err != nil {
				panic(err)
			}
		}()
	}

	handleOutput := func(shortID string, color string, timestamp int64, message string) {
		if runtime.GOOS == "windows" {
			fmt.Printf(
				"%v@%v %v\n",
				shortID,
				time.Now().Unix(),
				strings.TrimFunc(message, func(r rune) bool {
					return !unicode.IsGraphic(r)
				}),
			)
		} else {
			fmt.Printf(
				"%v%v%v@%v%v %v%v%v\n",
				utils.ColorBackgroundBlack,
				color,
				shortID,
				time.Now().Unix(),
				utils.ColorReset,
				color,
				strings.TrimFunc(message, func(r rune) bool {
					return !unicode.IsGraphic(r)
				}),
				utils.ColorReset,
			)
		}
	}

	bdrs := []builders.Builder{
		apk.NewBuilder(
			ctx,
			cli,

			apk.Image,
			*pull,
			*src,
			filepath.Join(*dst, "apk"),
			handleID,
			handleOutput,
			*appID,
			*gpgKeyContent,
			*gpgKeyPassword,
			*apkCertContent,
			*apkCertPassword,
			*baseURL+"apk",
			*appID,
			false,
		),
		deb.NewBuilder(
			ctx,
			cli,

			deb.Image,
			*pull,
			*src,
			filepath.Join(*dst, "deb", "debian", "sid", "x86_64"),
			handleID,
			handleOutput,
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
			*src,
			filepath.Join(*dst, "dmg"),
			handleID,
			handleOutput,
			*appID,
			*appName,
			*gpgKeyContent,
			*gpgKeyPassword,
			*dmgUniversal,
			[]string{},
		),
		flatpak.NewBuilder(
			ctx,
			cli,

			flatpak.Image,
			*pull,
			*src,
			filepath.Join(*dst, "flatpak", "x86_64"),
			handleID,
			handleOutput,
			*appID,
			*gpgKeyContent,
			*gpgKeyPassword,
			*gpgKeyID,
			*baseURL+"flatpak/x86_64",
			"amd64",
		),
		msi.NewBuilder(
			ctx,
			cli,

			msi.Image,
			*pull,
			*src,
			filepath.Join(*dst, "msi", "x86_64"),
			handleID,
			handleOutput,
			*appID,
			*appName,
			*gpgKeyContent,
			*gpgKeyPassword,
			"amd64",
			[]string{},
		),
		rpm.NewBuilder(
			ctx,
			cli,

			rpm.Image,
			*pull,
			*src,
			filepath.Join(*dst, "rpm", "fedora", "36", "x86_64"),
			handleID,
			handleOutput,
			*appID,
			*gpgKeyContent,
			*gpgKeyPassword,
			*gpgKeyID,
			*baseURL,
			*rpmPackageVersion,
			"fedora-36",
			"amd64",
			*rpmPackageSuffix,
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

		break
	}

	wg.Wait()
}
