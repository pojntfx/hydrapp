package cmd

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/builders"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/builders/apk"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/builders/binaries"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/builders/deb"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/builders/dmg"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/builders/docs"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/builders/flatpak"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/builders/msi"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/builders/rpm"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/builders/tests"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/config"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/secrets"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	configFlag      = "config"
	pullFlag        = "pull"
	tagFlag         = "tag"
	concurrencyFlag = "concurrency"
	ejectFlag       = "eject"
	overwriteFlag   = "overwrite"
	srcFlag         = "src"
	dstFlag         = "dst"
	excludeFlag     = "exclude"

	javaKeystoreFlag = "java-keystore"

	pgpKeyFlag   = "pgp-key"
	pgpKeyIDFlag = "pgp-key-id"

	branchIDFlag        = "branch-id"
	branchNameFlag      = "branch-name"
	branchTimestampFlag = "branch-timestamp"
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

		configFile, err := os.Open(viper.GetString(configFlag))
		if err != nil {
			return err
		}
		defer configFile.Close()

		cfg, err := config.Parse(configFile)
		if err != nil {
			return err
		}

		var (
			branchID        = viper.GetString(branchIDFlag)
			branchName      = viper.GetString(branchNameFlag)
			branchTimestamp = time.Unix(viper.GetInt64(branchTimestampFlag), 0)
		)
		if !(viper.IsSet(branchIDFlag) && viper.IsSet(branchNameFlag) && viper.IsSet(branchTimestampFlag)) {
			repo, err := git.PlainOpen(viper.GetString(srcFlag))
			if err != nil && !errors.Is(err, git.ErrRepositoryNotExists) { // If source directory is not a Git repository, use provided flags
				return err
			} else if err == nil {
				headRef, err := repo.Head()
				if err != nil {
					return err
				}

				headCommit, err := repo.CommitObject(headRef.Hash())
				if err != nil {
					return err
				}

				tags, err := repo.Tags()
				if err != nil {
					return err
				}

				isTag := false
				if err := tags.ForEach(func(r *plumbing.Reference) error {
					if r.Hash() == headCommit.Hash {
						isTag = true
					}

					return nil
				}); err != nil {
					return err
				}

				if isTag {
					if !viper.IsSet(branchIDFlag) {
						branchID = ""
					}

					if !viper.IsSet(branchNameFlag) {
						branchName = ""
					}
				} else {
					if !viper.IsSet(branchIDFlag) {
						branchID = headRef.Name().Short()
					}

					if !viper.IsSet(branchNameFlag) {
						branchName = utils.Capitalize(branchID)
					}
				}

				if !viper.IsSet(branchTimestampFlag) {
					branchTimestamp = headCommit.Author.When
				}
			}
		}

		var (
			javaKeystore            []byte
			javaKeystorePassword    string
			javaCertificatePassword string

			pgpKey         []byte
			pgpKeyPassword string
			pgpKeyID       string
		)
		if !viper.GetBool(ejectFlag) {
			javaKeystorePassword = viper.GetString(javaKeystorePasswordFlag)
			javaCertificatePassword = viper.GetString(javaCertificatePasswordFlag)

			pgpKeyPassword = viper.GetString(pgpKeyPasswordFlag)
			pgpKeyID = viper.GetString(pgpKeyIDFlag)

			var scs secrets.Root
			if strings.TrimSpace(viper.GetString(javaKeystoreFlag)) == "" &&
				strings.TrimSpace(javaKeystorePassword) == "" &&
				strings.TrimSpace(javaCertificatePassword) == "" &&

				strings.TrimSpace(viper.GetString(pgpKeyFlag)) == "" &&
				strings.TrimSpace(pgpKeyPassword) == "" &&
				strings.TrimSpace(pgpKeyID) == "" {
				secretsFile, err := os.Open(viper.GetString(secretsFlag))
				if err == nil {
					defer secretsFile.Close()

					s, err := secrets.Parse(secretsFile)
					if err != nil {
						return err
					}
					scs = *s
				} else {
					if !errors.Is(err, os.ErrNotExist) {
						return err
					}

					keystorePassword, err := secrets.GeneratePassword(32)
					if err != nil {
						return err
					}

					certificatePassword, err := secrets.GeneratePassword(32)
					if err != nil {
						return err
					}

					keystoreBuf := &bytes.Buffer{}
					if err := secrets.GenerateKeystore(
						keystorePassword,
						certificatePassword,
						fullNameDefault,
						fullNameDefault,
						certificateValidityDefault,
						javaRSABitsDefault,
						keystoreBuf,
					); err != nil {
						return err
					}

					pgpPassword, err := secrets.GeneratePassword(32)
					if err != nil {
						return err
					}

					pgpKey, pgpKeyID, err := secrets.GeneratePGPKey(
						fullNameDefault,
						emailDefault,
						pgpPassword,
					)
					if err != nil {
						return err
					}

					scs = secrets.Root{
						JavaSecrets: secrets.JavaSecrets{
							Keystore:            keystoreBuf.Bytes(),
							KeystorePassword:    keystorePassword,
							CertificatePassword: certificatePassword,
						},
						PGPSecrets: secrets.PGPSecrets{
							Key:         pgpKey,
							KeyID:       pgpKeyID,
							KeyPassword: pgpPassword,
						},
					}

					if err := os.MkdirAll(filepath.Dir(viper.GetString(secretsFlag)), os.ModePerm); err != nil {
						return err
					}

					out, err := os.OpenFile(viper.GetString(secretsFlag), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
					if err != nil {
						return err
					}
					defer out.Close()

					if err := yaml.NewEncoder(out).Encode(scs); err != nil {
						return err
					}
				}
			}

			if strings.TrimSpace(viper.GetString(javaKeystoreFlag)) == "" {
				javaKeystore = scs.JavaSecrets.Keystore
			} else {
				javaKeystore, err = os.ReadFile(viper.GetString(javaKeystoreFlag))
				if err != nil {
					return err
				}
			}

			if strings.TrimSpace(javaKeystorePassword) == "" {
				javaKeystorePassword = base64.StdEncoding.EncodeToString([]byte(scs.JavaSecrets.KeystorePassword))
			}

			if strings.TrimSpace(javaCertificatePassword) == "" {
				javaCertificatePassword = base64.StdEncoding.EncodeToString([]byte(scs.JavaSecrets.CertificatePassword))
			}

			if strings.TrimSpace(viper.GetString(pgpKeyFlag)) == "" {
				pgpKey = []byte(scs.PGPSecrets.Key)
			} else {
				pgpKey, err = os.ReadFile(viper.GetString(pgpKeyFlag))
				if err != nil {
					return err
				}
			}

			if strings.TrimSpace(pgpKeyPassword) == "" {
				pgpKeyPassword = base64.StdEncoding.EncodeToString([]byte(scs.PGPSecrets.KeyPassword))
			}

			if strings.TrimSpace(pgpKeyID) == "" {
				pgpKeyID = base64.StdEncoding.EncodeToString([]byte(scs.PGPSecrets.KeyID))
			}
		}

		licenseText, err := os.ReadFile(filepath.Join(filepath.Dir(viper.GetString(configFlag)), "LICENSE"))
		if err != nil {
			return err
		}

		cli, err := client.NewClientWithOpts(client.FromEnv)
		if err != nil {
			return err
		}
		defer cli.Close()

		// See https://github.com/rancher/rke/issues/1711#issuecomment-578382159
		cli.NegotiateAPIVersion(ctx)

		handleID := func(id string) {
			s := make(chan os.Signal, 1)
			signal.Notify(s, os.Interrupt, syscall.SIGTERM)

			go func() {
				<-s

				log.Println("Gracefully shutting down")

				go func() {
					<-s

					log.Println("Forcing shutdown")

					os.Exit(1)
				}()

				if err := cli.ContainerRemove(ctx, id, container.RemoveOptions{
					Force: true,
				}); err != nil {
					panic(err)
				}
			}()
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
					os.Stdout,
					"icon.png",
					cfg.App.ID,
					pgpKey,
					pgpKeyPassword,
					pgpKeyID,
					cfg.App.BaseURL+"/"+c.Path,
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
					branchID,
					branchName,
					branchTimestamp,
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
						os.Stdout,
						"icon.png",
						cfg.App.ID,
						cfg.App.Name,
						pgpKey,
						pgpKeyPassword,
						cfg.DMG.Packages,
						cfg.Releases,
						viper.GetBool(overwriteFlag),
						branchID,
						branchName,
						branchTimestamp,
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
					os.Stdout,
					"icon.png",
					cfg.App.ID,
					pgpKey,
					pgpKeyPassword,
					pgpKeyID,
					cfg.App.BaseURL+"/"+c.Path,
					c.Architecture,
					cfg.App.Name,
					cfg.App.Description,
					cfg.App.Summary,
					cfg.App.License,
					cfg.App.Homepage,
					cfg.App.Git,
					c.Packages,
					cfg.Releases,
					viper.GetBool(overwriteFlag),
					branchID,
					branchName,
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
					os.Stdout,
					"icon.png",
					cfg.App.ID,
					cfg.App.Name,
					pgpKey,
					pgpKeyPassword,
					c.Architecture,
					c.Packages,
					cfg.Releases,
					viper.GetBool(overwriteFlag),
					branchID,
					branchName,
					branchTimestamp,
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
					os.Stdout,
					"icon.png",
					cfg.App.ID,
					pgpKey,
					pgpKeyPassword,
					pgpKeyID,
					cfg.App.BaseURL+"/"+c.Path,
					c.Distro,
					c.Architecture,
					c.Trailer,
					cfg.App.Name,
					cfg.App.Description,
					cfg.App.Summary,
					cfg.App.Homepage,
					cfg.App.Git,
					cfg.App.License,
					cfg.Releases,
					c.Packages,
					viper.GetBool(overwriteFlag),
					branchID,
					branchName,
					branchTimestamp,
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
						os.Stdout,
						cfg.App.ID,
						javaKeystore,
						javaKeystorePassword,
						javaCertificatePassword,
						pgpKey,
						pgpKeyPassword,
						cfg.App.BaseURL+"/"+cfg.APK.Path,
						cfg.App.Name,
						cfg.Releases,
						viper.GetBool(overwriteFlag),
						branchID,
						branchName,
						branchTimestamp,
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
						os.Stdout,
						cfg.App.ID,
						pgpKey,
						pgpKeyPassword,
						cfg.App.Name,
						branchID,
						branchName,
						branchTimestamp,
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
						os.Stdout,
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
						os.Stdout,
						branchID,
						branchName,
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

	buildCmd.PersistentFlags().String(configFlag, "hydrapp.yaml", "Config file to use")

	buildCmd.PersistentFlags().Bool(pullFlag, false, "Whether to (re-)pull the images or not")
	buildCmd.PersistentFlags().String(tagFlag, "latest", "Image tag to use")
	buildCmd.PersistentFlags().Int(concurrencyFlag, 1, "Maximum amount of concurrent builders to run at once")
	buildCmd.PersistentFlags().Bool(ejectFlag, false, "Write platform-specific config files (AndroidManifest.xml, .spec etc.) to directory specified by --src, then exit (--exclude still applies)")
	buildCmd.PersistentFlags().Bool(overwriteFlag, false, "Overwrite platform-specific config files even if they exist")

	buildCmd.PersistentFlags().String(srcFlag, pwd, "Source directory (must be absolute path)")
	buildCmd.PersistentFlags().String(dstFlag, filepath.Join(pwd, "out"), "Output directory (must be absolute path)")

	buildCmd.PersistentFlags().String(excludeFlag, "", "Regex of platforms and architectures not to build for, i.e. (binaries|deb|rpm|flatpak/amd64|msi/386|dmg|docs|tests)")

	buildCmd.PersistentFlags().String(javaKeystoreFlag, "", "Path to Java/APK keystore (neither path nor content should be not base64-encoded)")
	buildCmd.PersistentFlags().String(javaKeystorePasswordFlag, "", "Java/APK keystore password (base64-encoded)")
	buildCmd.PersistentFlags().String(javaCertificatePasswordFlag, "", " Java/APK certificate password (base64-encoded) (if keystore uses PKCS12, this will be the same as --java-keystore-password)")

	buildCmd.PersistentFlags().String(pgpKeyFlag, "", "Path to armored PGP private key (neither path nor content should be not base64-encoded)")
	buildCmd.PersistentFlags().String(pgpKeyPasswordFlag, "", "PGP key password (base64-encoded)")
	buildCmd.PersistentFlags().String(pgpKeyIDFlag, "", "PGP key ID (base64-encoded)")

	buildCmd.PersistentFlags().String(branchIDFlag, "", `Branch ID to build the app as, i.e. main (for an app ID like "myappid.main" and baseURL like "mybaseurl/main") (fetched from Git unless set)`)
	buildCmd.PersistentFlags().String(branchNameFlag, "", `Branch name to build the app as, i.e. Main (for an app name like "myappname (Main)") (fetched from Git unless set)`)
	buildCmd.PersistentFlags().Int64(branchTimestampFlag, 0, `Branch UNIX timestamp to build the app with, i.e. 1715484587 (fetched from Git unless set)`)

	viper.AutomaticEnv()

	rootCmd.AddCommand(buildCmd)
}
