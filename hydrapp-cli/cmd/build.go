package cmd

import (
	"context"
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
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/builders"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/builders/apk"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/builders/binaries"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/builders/deb"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/builders/dmg"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/builders/docs"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/builders/flatpak"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/builders/msi"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/builders/rpm"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/builders/tests"
	"github.com/pojntfx/hydrapp/hydrapp-cli/pkg/utils"
	cconfig "github.com/pojntfx/hydrapp/hydrapp-utils/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configFlag       = "config"
	pullFlag         = "pull"
	tagFlag          = "tag"
	concurrencyFlag  = "concurrency"
	ejectFlag        = "eject"
	overwriteFlag    = "overwrite"
	srcFlag          = "src"
	dstFlag          = "dst"
	excludeFlag      = "exclude"
	pgpKeyFlag       = "pgp-key"
	pgpPasswordFlag  = "pgp-password"
	pgpIDFlag        = "pgp-id"
	apkCertFlag      = "apk-cert"
	apkStorepassFlag = "apk-storepass"
	apkKeypassFlag   = "apk-keypass"
	branchIDFlag     = "branch-id"
	branchNameFlag   = "branch-name"
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

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b"},
	Short:   "Build a hydrapp project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		content, err := ioutil.ReadFile(viper.GetString(configFlag))
		if err != nil {
			return err
		}

		cfg, err := cconfig.Parse(content)
		if err != nil {
			return err
		}

		licenseText, err := ioutil.ReadFile(filepath.Join(filepath.Dir(viper.GetString(configFlag)), "LICENSE"))
		if err != nil {
			return err
		}

		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return err
		}
		defer cli.Close()

		var pgpKeyContent []byte
		if !viper.GetBool(ejectFlag) {
			pgpKeyContent, err = ioutil.ReadFile(viper.GetString(pgpKeyFlag))
			if err != nil {
				return err
			}
		}

		var apkCertContent []byte
		if !viper.GetBool(ejectFlag) {
			apkCertContent, err = ioutil.ReadFile(viper.GetString(apkCertFlag))
			if err != nil {
				return err
			}
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
			skip, err := checkIfSkip(viper.GetString(excludeFlag), "deb", c.Architecture)
			if err != nil {
				return err
			}

			if skip {
				continue
			}

			bdrs = append(
				bdrs,
				deb.NewBuilder(
					ctx,
					cli,

					deb.Image+":"+viper.GetString(tagFlag),
					viper.GetBool(pullFlag),
					viper.GetString(srcFlag),
					filepath.Join(viper.GetString(dstFlag), c.Path),
					handleID,
					handleOutput,
					cfg.App.ID,
					pgpKeyContent,
					viper.GetString(pgpPasswordFlag),
					viper.GetString(pgpIDFlag),
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
					cfg.App.License,
					string(licenseText),
					cfg.App.Name,
					viper.GetBool(overwriteFlag),
					viper.GetString(branchIDFlag),
					viper.GetString(branchNameFlag),
					cfg.Go.Main,
					cfg.Go.Flags,
					cfg.Go.Generate,
				),
			)
		}

		if strings.TrimSpace(cfg.DMG.Path) != "" {
			skip, err := checkIfSkip(viper.GetString(excludeFlag), "dmg", "")
			if err != nil {
				return err
			}

			if !skip {
				bdrs = append(
					bdrs,
					dmg.NewBuilder(
						ctx,
						cli,

						dmg.Image+":"+viper.GetString(tagFlag),
						viper.GetBool(pullFlag),
						viper.GetString(srcFlag),
						filepath.Join(viper.GetString(dstFlag), cfg.DMG.Path),
						handleID,
						handleOutput,
						cfg.App.ID,
						cfg.App.Name,
						pgpKeyContent,
						viper.GetString(pgpPasswordFlag),
						cfg.DMG.Packages,
						cfg.Releases,
						viper.GetBool(overwriteFlag),
						viper.GetString(branchIDFlag),
						viper.GetString(branchNameFlag),
						cfg.Go.Main,
						cfg.Go.Flags,
						cfg.Go.Generate,
					),
				)
			}
		}

		for _, c := range cfg.Flatpak {
			skip, err := checkIfSkip(viper.GetString(excludeFlag), "flatpak", c.Architecture)
			if err != nil {
				return err
			}

			if skip {
				continue
			}

			bdrs = append(
				bdrs,
				flatpak.NewBuilder(
					ctx,
					cli,

					flatpak.Image+":"+viper.GetString(tagFlag),
					viper.GetBool(pullFlag),
					viper.GetString(srcFlag),
					filepath.Join(viper.GetString(dstFlag), c.Path),
					handleID,
					handleOutput,
					cfg.App.ID,
					pgpKeyContent,
					viper.GetString(pgpPasswordFlag),
					viper.GetString(pgpIDFlag),
					cfg.App.BaseURL+c.Path,
					c.Architecture,
					cfg.App.Name,
					cfg.App.Description,
					cfg.App.Summary,
					cfg.App.License,
					cfg.App.Homepage,
					cfg.Releases,
					viper.GetBool(overwriteFlag),
					viper.GetString(branchIDFlag),
					viper.GetString(branchNameFlag),
					cfg.Go.Main,
					cfg.Go.Flags,
					cfg.Go.Generate,
				),
			)
		}

		for _, c := range cfg.MSI {
			skip, err := checkIfSkip(viper.GetString(excludeFlag), "msi", c.Architecture)
			if err != nil {
				return err
			}

			if skip {
				continue
			}

			bdrs = append(
				bdrs,
				msi.NewBuilder(
					ctx,
					cli,

					msi.Image+":"+viper.GetString(tagFlag),
					viper.GetBool(pullFlag),
					viper.GetString(srcFlag),
					filepath.Join(viper.GetString(dstFlag), c.Path),
					handleID,
					handleOutput,
					cfg.App.ID,
					cfg.App.Name,
					pgpKeyContent,
					viper.GetString(pgpPasswordFlag),
					c.Architecture,
					c.Packages,
					cfg.Releases,
					viper.GetBool(overwriteFlag),
					viper.GetString(branchIDFlag),
					viper.GetString(branchNameFlag),
					cfg.Go.Main,
					cfg.Go.Flags,
					c.Include,
					cfg.Go.Generate,
				),
			)
		}

		for _, c := range cfg.RPM {
			skip, err := checkIfSkip(viper.GetString(excludeFlag), "rpm", c.Architecture)
			if err != nil {
				return err
			}

			if skip {
				continue
			}

			bdrs = append(
				bdrs,
				rpm.NewBuilder(
					ctx,
					cli,

					rpm.Image+":"+viper.GetString(tagFlag),
					viper.GetBool(pullFlag),
					viper.GetString(srcFlag),
					filepath.Join(viper.GetString(dstFlag), c.Path),
					handleID,
					handleOutput,
					cfg.App.ID,
					pgpKeyContent,
					viper.GetString(pgpPasswordFlag),
					viper.GetString(pgpIDFlag),
					cfg.App.BaseURL+c.Path,
					c.Distro,
					c.Architecture,
					c.Trailer,
					cfg.App.Name,
					cfg.App.Description,
					cfg.App.Summary,
					cfg.App.Homepage,
					cfg.App.License,
					cfg.Releases,
					c.Packages,
					viper.GetBool(overwriteFlag),
					viper.GetString(branchIDFlag),
					viper.GetString(branchNameFlag),
					cfg.Go.Main,
					cfg.Go.Flags,
					cfg.Go.Generate,
				),
			)
		}

		if strings.TrimSpace(cfg.APK.Path) != "" {
			skip, err := checkIfSkip(viper.GetString(excludeFlag), "apk", "")
			if err != nil {
				return err
			}

			if !skip {
				bdrs = append(
					bdrs,
					apk.NewBuilder(
						ctx,
						cli,

						apk.Image+":"+viper.GetString(tagFlag),
						viper.GetBool(pullFlag),
						viper.GetString(srcFlag),
						filepath.Join(viper.GetString(dstFlag), cfg.APK.Path),
						handleID,
						handleOutput,
						cfg.App.ID,
						pgpKeyContent,
						viper.GetString(pgpPasswordFlag),
						apkCertContent,
						viper.GetString(apkStorepassFlag),
						viper.GetString(apkKeypassFlag),
						cfg.App.BaseURL+cfg.APK.Path,
						cfg.App.Name,
						viper.GetBool(overwriteFlag),
						viper.GetString(branchIDFlag),
						viper.GetString(branchNameFlag),
						cfg.Go.Main,
						cfg.Go.Flags,
						cfg.Go.Generate,
					),
				)
			}
		}

		if strings.TrimSpace(cfg.Binaries.Path) != "" {
			skip, err := checkIfSkip(viper.GetString(excludeFlag), "binaries", "")
			if err != nil {
				return err
			}

			if !skip {
				bdrs = append(
					bdrs,
					binaries.NewBuilder(
						ctx,
						cli,

						binaries.Image+":"+viper.GetString(tagFlag),
						viper.GetBool(pullFlag),
						viper.GetString(srcFlag),
						filepath.Join(viper.GetString(dstFlag), cfg.Binaries.Path),
						handleID,
						handleOutput,
						cfg.App.ID,
						pgpKeyContent,
						viper.GetString(pgpPasswordFlag),
						cfg.App.Name,
						viper.GetString(branchIDFlag),
						viper.GetString(branchNameFlag),
						cfg.Go.Main,
						cfg.Go.Flags,
						cfg.Go.Generate,
						cfg.Binaries.Exclude,
						cfg.Binaries.Packages,
					),
				)
			}
		}

		if strings.TrimSpace(cfg.Go.Tests) != "" {
			skip, err := checkIfSkip(viper.GetString(excludeFlag), "tests", "")
			if err != nil {
				return err
			}

			if !skip {
				bdrs = append(
					bdrs,
					tests.NewBuilder(
						ctx,
						cli,

						cfg.Go.Image,
						viper.GetBool(pullFlag),
						viper.GetString(srcFlag),
						"",
						handleID,
						handleOutput,
						cfg.Go.Flags,
						cfg.Go.Generate,
						cfg.Go.Tests,
					),
				)
			}
		}

		if strings.TrimSpace(cfg.Docs.Path) != "" {
			skip, err := checkIfSkip(viper.GetString(excludeFlag), "docs", "")
			if err != nil {
				return err
			}

			if !skip {
				bdrs = append(
					bdrs,
					docs.NewBuilder(
						ctx,
						cli,

						docs.Image+":"+viper.GetString(tagFlag),
						viper.GetBool(pullFlag),
						viper.GetString(srcFlag),
						filepath.Join(viper.GetString(dstFlag), cfg.Docs.Path),
						handleID,
						handleOutput,
						viper.GetString(branchIDFlag),
						viper.GetString(branchNameFlag),
						cfg.Go.Main,
						cfg,
						viper.GetBool(overwriteFlag),
					),
				)
			}
		}

		semaphore := make(chan struct{}, viper.GetInt(concurrencyFlag))
		var wg sync.WaitGroup
		for _, b := range bdrs {
			wg.Add(1)

			semaphore <- struct{}{}

			go func(builder builders.Builder) {
				defer func() {
					<-semaphore

					wg.Done()
				}()

				if viper.GetBool(ejectFlag) {
					if err := builder.Render(viper.GetString(srcFlag), true); err != nil {
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

		return nil
	},
}

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	buildCmd.PersistentFlags().Bool(noNetworkFlag, false, "Disable all network interaction")

	buildCmd.PersistentFlags().String(configFlag, "hydrapp.yaml", "Config file to use")

	buildCmd.PersistentFlags().Bool(pullFlag, false, "Whether to pull the images or not")
	buildCmd.PersistentFlags().String(tagFlag, "latest", "Image tag to use")
	buildCmd.PersistentFlags().Int(concurrencyFlag, 1, "Maximum amount of concurrent builders to run at once")
	buildCmd.PersistentFlags().Bool(ejectFlag, false, "Write platform-specific config files (AndroidManifest.xml, .spec etc.) to directory specified by --src, then exit (--exclude still applies")
	buildCmd.PersistentFlags().Bool(overwriteFlag, false, "Overwrite platform-specific config files even if they exist")

	buildCmd.PersistentFlags().String(srcFlag, pwd, "Source directory (must be absolute path)")
	buildCmd.PersistentFlags().String(dstFlag, filepath.Join(pwd, "out"), "Output directory (must be absolute path)")

	buildCmd.PersistentFlags().String(excludeFlag, "", "Regex of platforms and architectures not to build for, i.e. (apk|dmg|msi/386|flatpak/amd64)")

	buildCmd.PersistentFlags().String(pgpKeyFlag, "", "Path to armored PGP private key")
	buildCmd.PersistentFlags().String(pgpPasswordFlag, "", "Password for PGP key")
	buildCmd.PersistentFlags().String(pgpIDFlag, "", "ID of the PGP key to use")

	buildCmd.PersistentFlags().String(apkCertFlag, "", "Path to Android keystore")
	buildCmd.PersistentFlags().String(apkStorepassFlag, "", "Password for Android keystore")
	buildCmd.PersistentFlags().String(apkKeypassFlag, "", " Password for Android certificate (if keystore uses PKCS12, this will be the same as --apk-storepass)")

	buildCmd.PersistentFlags().String(branchIDFlag, "", `Branch ID to build the app as, i.e. main (for an app ID like "myappid.main" and baseURL like "mybaseurl/main"`)
	buildCmd.PersistentFlags().String(branchNameFlag, "", `Branch name to build the app as, i.e. Main (for an app name like "myappname (Main)"`)

	viper.AutomaticEnv()

	rootCmd.AddCommand(buildCmd)
}
