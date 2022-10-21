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
	"regexp"
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

func checkIfSkip(exclude string, platform, architecture string) (bool, error) {
	if strings.TrimSpace(exclude) == "" {
		return false, nil
	}

	skip, err := regexp.MatchString(exclude, platform+"/"+architecture)
	if err != nil {
		return false, err
	}

	if skip {
		log.Printf("Skipping %v/%v (platform or architecture matched the provided regex)", platform, architecture)

		return true, nil
	}

	return false, nil
}

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	config := flag.String("config", "hydrapp.yaml", "Config file to use")

	pull := flag.Bool("pull", false, "Whether to pull the images or not")
	concurrency := flag.Int("concurrency", 1, "Maximum amount of concurrent builders to run at once")
	eject := flag.Bool("eject", false, "Write platform-specific config files (AndroidManifest.xml, .spec etc.) to directory specified by --src, then exit (--exclude still applies")
	overwrite := flag.Bool("overwrite", false, "Overwrite platform-specific config files even if they exist")

	src := flag.String("src", pwd, "Source directory (must be absolute path)")
	dst := flag.String("dst", filepath.Join(pwd, "out"), "Output directory (must be absolute path)")

	exclude := flag.String("exclude", "", "Regex of platforms and architectures not to build for, i.e. (apk|dmg|msi/386|flatpak/amd64)")

	gpgKey := flag.String("gpg-key", "", "Path to armored GPG private key")
	gpgPassword := flag.String("gpg-password", "", "Password for GPG key")
	gpgID := flag.String("gpg-id", "", "ID of the GPG key to use")

	apkCert := flag.String("apk-cert", "", "Path to Android certificate/keystore")
	apkPassword := flag.String("apk-password", "", " Password for Android certificate")

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

	gpgKeyContent, err := ioutil.ReadFile(*gpgKey)
	if err != nil {
		panic(err)
	}

	apkCertContent, err := ioutil.ReadFile(*apkCert)
	if err != nil {
		panic(err)
	}

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
		skip, err := checkIfSkip(*exclude, "deb", c.Architecture)
		if err != nil {
			panic(err)
		}

		if skip {
			continue
		}

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
				gpgKeyContent,
				*gpgPassword,
				*gpgID,
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
				*overwrite,
			),
		)
	}

	if strings.TrimSpace(cfg.DMG.Path) != "" {
		skip, err := checkIfSkip(*exclude, "dmg", "")
		if err != nil {
			panic(err)
		}

		if !skip {
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
					gpgKeyContent,
					*gpgPassword,
					cfg.DMG.Universal,
					cfg.DMG.Packages,
					cfg.Releases,
					*overwrite,
				),
			)
		}
	}

	for _, c := range cfg.Flatpak {
		skip, err := checkIfSkip(*exclude, "flatpak", c.Architecture)
		if err != nil {
			panic(err)
		}

		if skip {
			continue
		}

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
				gpgKeyContent,
				*gpgPassword,
				*gpgID,
				cfg.App.BaseURL+c.Path,
				c.Architecture,
				cfg.App.Name,
				cfg.App.Description,
				cfg.App.Summary,
				cfg.License.SPDX,
				cfg.App.Homepage,
				cfg.Releases,
				*overwrite,
			),
		)
	}

	for _, c := range cfg.MSI {
		skip, err := checkIfSkip(*exclude, "msi", c.Architecture)
		if err != nil {
			panic(err)
		}

		if skip {
			continue
		}

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
				gpgKeyContent,
				*gpgPassword,
				c.Architecture,
				c.Packages,
				cfg.Releases,
				*overwrite,
			),
		)
	}

	for _, c := range cfg.RPM {
		skip, err := checkIfSkip(*exclude, "rpm", c.Architecture)
		if err != nil {
			panic(err)
		}

		if skip {
			continue
		}

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
				gpgKeyContent,
				*gpgPassword,
				*gpgID,
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
				*overwrite,
			),
		)
	}

	if strings.TrimSpace(cfg.APK.Path) != "" {
		skip, err := checkIfSkip(*exclude, "apk", "")
		if err != nil {
			panic(err)
		}

		if !skip {
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
					gpgKeyContent,
					*gpgPassword,
					apkCertContent,
					*apkPassword,
					cfg.App.BaseURL+cfg.APK.Path,
					cfg.App.ID,
					*overwrite,
				),
			)
		}
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

			if *eject {
				if err := builder.Render(*src); err != nil {
					panic(err)
				}
			} else {
				if err := builder.Build(); err != nil {
					panic(err)
				}
			}
		}(b)
	}

	wg.Wait()
}
