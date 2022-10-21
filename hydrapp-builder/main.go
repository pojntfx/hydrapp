package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
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
	cconfig "github.com/pojntfx/hydrapp/hydrapp-builder/pkg/config"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/utils"
)

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	config := flag.String("config", "hydrapp.yaml", "Config file to use")

	pull := flag.Bool("pull", false, "Whether to pull the images or not")
	concurrency := flag.Int("concurrency", 1, "Maximum amount of concurrent builders to run at once")

	src := flag.String("src", pwd, "Source directory")
	dst := flag.String("dst", filepath.Join(pwd, "out"), "Output directory")

	gpgKeyContent := flag.String("gpg-key-content", "", "base64-encoded GPG key contents")
	gpgKeyPassword := flag.String("gpg-key-password", "", " base64-encoded password for the GPG key")
	gpgKeyID := flag.String("gpg-key-id", "", "ID of the GPG key to use")

	apkCertContent := flag.String("apk-cert-content", "", "base64-encoded Android cert contents")
	apkCertPassword := flag.String("apk-cert-password", "", " base64-encoded password for the Android cert")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	content, err := ioutil.ReadFile(*config)
	if err != nil {
		panic(err)
	}

	cfg, err := cconfig.Parse(content)
	if err != nil {
		panic(err)
	}

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

	bdrs := []builders.Builder{}

	for _, c := range cfg.DEB {
		bdrs = append(
			bdrs,
			deb.NewBuilder(
				ctx,
				cli,

				deb.Image,
				*pull,
				*src,
				filepath.Join(*dst, c.Path),
				handleID,
				handleOutput,
				cfg.App.ID,
				*gpgKeyContent,
				*gpgKeyPassword,
				*gpgKeyID,
				cfg.App.BaseURL+c.Path,
				c.OS,
				c.Distro,
				c.Mirrorsite,
				c.Components,
				c.Debootstrapopts,
				c.Architecture,
				cfg.Releases,
				cfg.App.Description,
				cfg.App.Summary,
				cfg.App.Homepage,
				cfg.App.Git,
				c.Packages,
				cfg.License.SPDX,
				cfg.License.Text,
				cfg.App.Name,
				false,
			),
		)
	}

	if strings.TrimSpace(cfg.DMG.Path) != "" {
		bdrs = append(
			bdrs,
			dmg.NewBuilder(
				ctx,
				cli,

				dmg.Image,
				*pull,
				*src,
				filepath.Join(*dst, cfg.DMG.Path),
				handleID,
				handleOutput,
				cfg.App.ID,
				cfg.App.Name,
				*gpgKeyContent,
				*gpgKeyPassword,
				cfg.DMG.Universal,
				cfg.DMG.Packages,
				cfg.Releases,
				false,
			),
		)
	}

	for _, c := range cfg.Flatpak {
		bdrs = append(
			bdrs,
			flatpak.NewBuilder(
				ctx,
				cli,

				flatpak.Image,
				*pull,
				*src,
				filepath.Join(*dst, c.Path),
				handleID,
				handleOutput,
				cfg.App.ID,
				*gpgKeyContent,
				*gpgKeyPassword,
				*gpgKeyID,
				cfg.App.BaseURL+c.Path,
				c.Architecture,
				cfg.App.Name,
				cfg.App.Description,
				cfg.App.Summary,
				cfg.License.SPDX,
				cfg.App.Homepage,
				cfg.Releases,
				false,
			),
		)
	}

	for _, c := range cfg.MSI {
		bdrs = append(
			bdrs,
			msi.NewBuilder(
				ctx,
				cli,

				msi.Image,
				*pull,
				*src,
				filepath.Join(*dst, c.Path),
				handleID,
				handleOutput,
				cfg.App.ID,
				cfg.App.Name,
				*gpgKeyContent,
				*gpgKeyPassword,
				c.Architecture,
				c.Packages,
				cfg.Releases,
				false,
			),
		)
	}

	for _, c := range cfg.RPM {
		bdrs = append(
			bdrs,
			rpm.NewBuilder(
				ctx,
				cli,

				rpm.Image,
				*pull,
				*src,
				filepath.Join(*dst, c.Path),
				handleID,
				handleOutput,
				cfg.App.ID,
				*gpgKeyContent,
				*gpgKeyPassword,
				*gpgKeyID,
				cfg.App.BaseURL,
				c.Distro,
				c.Architecture,
				c.Trailer,
				cfg.App.Name,
				cfg.App.Description,
				cfg.App.Summary,
				cfg.App.Homepage,
				cfg.License.SPDX,
				cfg.Releases,
				c.Packages,
				false,
			),
		)
	}

	if strings.TrimSpace(cfg.APK.Path) != "" {
		bdrs = append(
			bdrs,
			apk.NewBuilder(
				ctx,
				cli,

				apk.Image,
				*pull,
				*src,
				filepath.Join(*dst, cfg.APK.Path),
				handleID,
				handleOutput,
				cfg.App.ID,
				*gpgKeyContent,
				*gpgKeyPassword,
				*apkCertContent,
				*apkCertPassword,
				cfg.App.BaseURL+cfg.APK.Path,
				cfg.App.ID,
				false,
			),
		)
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
