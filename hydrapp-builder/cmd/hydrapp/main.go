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

	pgpKey := flag.String("pgp-key", "", "Path to armored PGP private key")
	pgpPassword := flag.String("pgp-password", "", "Password for PGP key")
	pgpID := flag.String("pgp-id", "", "ID of the PGP key to use")

	apkCert := flag.String("apk-cert", "", "Path to Android keystore")
	apkStorepass := flag.String("apk-storepass", "", "Password for Android keystore")
	apkKeypass := flag.String("apk-keypass", "", " Password for Android certificate (if keystore uses PKCS12, this will be the same as --apk-storepass)")

	branchID := flag.String("branch-id", "", `Branch ID to build the app as, i.e. unstable (for an app ID like "myappid.unstable" and baseURL like "mybaseurl/unstable"`)
	branchName := flag.String("branch-name", "", `Branch name to build the app as, i.e. Unstable (for an app name like "myappname (Unstable)"`)

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

	pgpKeyContent, err := ioutil.ReadFile(*pgpKey)
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
				pgpKeyContent,
				*pgpPassword,
				*pgpID,
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
				*branchID,
				*branchName,
				cfg.Go.Main,
				cfg.Go.Flags,
				cfg.Go.Generate,
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
					pgpKeyContent,
					*pgpPassword,
					cfg.DMG.Universal,
					cfg.DMG.Packages,
					cfg.Releases,
					*overwrite,
					*branchID,
					*branchName,
					cfg.Go.Main,
					cfg.Go.Flags,
					cfg.Go.Generate,
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
				pgpKeyContent,
				*pgpPassword,
				*pgpID,
				cfg.App.BaseURL+c.Path,
				c.Architecture,
				cfg.App.Name,
				cfg.App.Description,
				cfg.App.Summary,
				cfg.License.SPDX,
				cfg.App.Homepage,
				cfg.Releases,
				*overwrite,
				*branchID,
				*branchName,
				cfg.Go.Main,
				cfg.Go.Flags,
				cfg.Go.Generate,
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
				pgpKeyContent,
				*pgpPassword,
				c.Architecture,
				c.Packages,
				cfg.Releases,
				*overwrite,
				*branchID,
				*branchName,
				cfg.Go.Main,
				cfg.Go.Flags,
				cfg.Go.Generate,
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
				pgpKeyContent,
				*pgpPassword,
				*pgpID,
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
				*branchID,
				*branchName,
				cfg.Go.Main,
				cfg.Go.Flags,
				cfg.Go.Generate,
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
					pgpKeyContent,
					*pgpPassword,
					apkCertContent,
					*apkStorepass,
					*apkKeypass,
					cfg.App.BaseURL+cfg.APK.Path,
					cfg.App.Name,
					*overwrite,
					*branchID,
					*branchName,
					cfg.Go.Main,
					cfg.Go.Flags,
					cfg.Go.Generate,
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
				if err := builder.Render(*src, true); err != nil {
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
