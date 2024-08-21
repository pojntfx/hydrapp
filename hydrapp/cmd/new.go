package cmd

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/config"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/generators"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/renderers/rpm"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	noNetworkFlag                     = "no-network"
	experimentalGithubPagesActionFlag = "experimental-github-pages-action"

	vanillaJSRESTKey  = "vanillajs-rest"
	vanillaJSFormsKey = "vanillajs-forms"
	reactPanrpcKey    = "react-panrpc"
)

var (
	errUnknownProjectType = errors.New("unknown project type")

	projectTypeItems = []generators.ProjectTypeOption{
		{
			Name:        vanillaJSRESTKey,
			Description: "Simple starter project with a REST API to connect the Vanilla JS frontend and backend",
		},
		{
			Name:        vanillaJSFormsKey,
			Description: "Traditional starter project with Web 1.0-style forms to connect the Vanilla JS frontend and backend",
		},
		{
			Name:        reactPanrpcKey,
			Description: "Complete starter project with panrpc RPCs to connect the React frontend and backend",
		},
	}
)

var newCmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"n"},
	Short:   "Generate a new hydrapp project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
			return err
		}

		projectTypeIndex, _, err := (&promptui.Select{
			Templates: &promptui.SelectTemplates{
				Label:    fmt.Sprintf("%s {{.}}: ", promptui.IconInitial),
				Active:   fmt.Sprintf("%s {{ .Name | underline }}: {{ .Description | faint }}", promptui.IconSelect),
				Inactive: "  {{ .Name }}: {{ .Description | faint }}",
				Selected: fmt.Sprintf(`{{ "%s" | green }} {{ .Name | faint }}: {{ .Description | faint }}`, promptui.IconGood),
			},
			Label: "Which project type do you want to generate?",
			Items: projectTypeItems,
		}).Run()
		if err != nil {
			return err
		}

		appID, err := (&promptui.Prompt{
			Label:   "App ID in reverse domain notation",
			Default: "com.github.example.myapp",
		}).Run()
		if err != nil {
			return err
		}

		appName, err := (&promptui.Prompt{
			Label:   "App name",
			Default: "My App",
		}).Run()
		if err != nil {
			return err
		}

		appSummary, err := (&promptui.Prompt{
			Label:   "App summary",
			Default: "My first app",
		}).Run()
		if err != nil {
			return err
		}

		appDescription, err := (&promptui.Prompt{
			Label:   "App description",
			Default: "My first application, built with hydrapp.",
		}).Run()
		if err != nil {
			return err
		}

		appHomepage, err := (&promptui.Prompt{
			Label:   "App homepage",
			Default: "https://github.com/example/myapp",
		}).Run()
		if err != nil {
			return err
		}

		appGit, err := (&promptui.Prompt{
			Label:   "App git repo",
			Default: appHomepage + ".git",
		}).Run()
		if err != nil {
			return err
		}

		appBaseurl, err := (&promptui.Prompt{
			Label:   "App base URL to expect the built assets to be published to",
			Default: "https://example.github.io/myapp/",
		}).Run()
		if err != nil {
			return err
		}

		goMod, err := (&promptui.Prompt{
			Label:   "Go module name",
			Default: "github.com/example/myapp",
		}).Run()
		if err != nil {
			return err
		}

		licenseSPDX, err := (&promptui.Prompt{
			Label:   "License SPDX identifier (see https://spdx.org/licenses/)",
			Default: "Apache-2.0",
		}).Run()
		if err != nil {
			return err
		}

		releaseAuthor, err := (&promptui.Prompt{
			Label:   "Release author name",
			Default: "Jean Doe",
		}).Run()
		if err != nil {
			return err
		}

		releaseEmail, err := (&promptui.Prompt{
			Label:   "Release author email",
			Default: "jean.doe@example.com",
		}).Run()
		if err != nil {
			return err
		}

		dir, err := (&promptui.Prompt{
			Label:   "Directory to write the app to",
			Default: "myapp",
		}).Run()
		if err != nil {
			return err
		}

		_, advancedConfiguration, err := (&promptui.Select{
			Label: "Do you want to do any advanced configuration?",
			Items: []string{"no", "yes"},
		}).Run()
		if err != nil {
			return err
		}

		goMain := "."
		goFlags := ""
		goGenerate := "go generate ./..."
		goTests := "go test ./..."
		goImg := "ghcr.io/pojntfx/hydrapp-build-tests:main"

		debArchitectures := "amd64"
		flatpakArchitectures := "amd64"
		msiArchitectures := "amd64"
		rpmArchitectures := "amd64"

		binariesExclude := "(android/*|ios/*|plan9/*|aix/*|linux/loong64|freebsd/riscv64|wasip1/wasm|js/wasm|openbsd/mips64)"

		if advancedConfiguration == "yes" {
			goMain, err = (&promptui.Prompt{
				Label:   "Go main package path",
				Default: goMain,
			}).Run()
			if err != nil {
				return err
			}

			goFlags, err = (&promptui.Prompt{
				Label:   "Go flags to pass to the compiler",
				Default: goFlags,
			}).Run()
			if err != nil {
				return err
			}

			goGenerate, err = (&promptui.Prompt{
				Label:   "Go generate command to run",
				Default: goGenerate,
			}).Run()
			if err != nil {
				return err
			}

			goTests, err = (&promptui.Prompt{
				Label:   "Go test command to run",
				Default: goTests,
			}).Run()
			if err != nil {
				return err
			}

			goImg, err = (&promptui.Prompt{
				Label:   "Go test OCI image to use",
				Default: goImg,
			}).Run()
			if err != nil {
				return err
			}

			debArchitectures, err = (&promptui.Prompt{
				Label:   "DEB architectures to build for (comma-seperated list of GOARCH values)",
				Default: debArchitectures,
			}).Run()
			if err != nil {
				return err
			}

			flatpakArchitectures, err = (&promptui.Prompt{
				Label:   "Flatpak architectures to build for (comma-seperated list of GOARCH values)",
				Default: flatpakArchitectures,
			}).Run()
			if err != nil {
				return err
			}

			msiArchitectures, err = (&promptui.Prompt{
				Label:   "MSI architectures to build for (comma-seperated list of GOARCH values)",
				Default: msiArchitectures,
			}).Run()
			if err != nil {
				return err
			}

			rpmArchitectures, err = (&promptui.Prompt{
				Label:   "RPM architectures to build for (comma-seperated list of GOARCH values)",
				Default: rpmArchitectures,
			}).Run()
			if err != nil {
				return err
			}

			binariesExclude, err = (&promptui.Prompt{
				Label:   "Regex of binaries to exclude from compilation",
				Default: binariesExclude,
			}).Run()
			if err != nil {
				return err
			}
		}

		licenseText := ""
		if !viper.GetBool(noNetworkFlag) {
			log.Println("Fetching full license text from SPDX ...")

			res, err := http.Get("https://raw.githubusercontent.com/spdx/license-list-data/main/text/" + licenseSPDX + ".txt")
			if err != nil {
				return err
			}
			if res.Body != nil {
				defer res.Body.Close()
			}
			if res.StatusCode != http.StatusOK {
				panic(errors.New(res.Status))
			}

			b, err := io.ReadAll(res.Body)
			if err != nil {
				return err
			}

			licenseText = string(b)

			log.Println("Success!")
		}

		{
			cfg := config.Root{}
			cfg.App = config.App{
				ID:          appID,
				Name:        appName,
				Summary:     appSummary,
				Description: appDescription,
				License:     licenseSPDX,
				Homepage:    appHomepage,
				Git:         appGit,
				BaseURL:     appBaseurl,
			}
			cfg.Go = config.Go{
				Main:     goMain,
				Flags:    goFlags,
				Generate: goGenerate,
				Tests:    goTests,
				Image:    goImg,
			}
			cfg.Releases = []renderers.Release{
				{
					Version:     "0.0.1",
					Date:        time.Now(),
					Description: "Initial release",
					Author:      releaseAuthor,
					Email:       releaseEmail,
				},
			}

			cfg.APK = config.APK{
				Path: "apk",
			}

			debs := []config.DEB{}
			for _, arch := range strings.Split(debArchitectures, ",") {
				debs = append(debs, config.DEB{
					Path:            path.Join("deb", "debian", "sid", utils.GetArchIdentifier(arch)),
					OS:              "debian",
					Distro:          "sid",
					Mirrorsite:      "http://http.us.debian.org/debian",
					Components:      []string{"main", "contrib"},
					Debootstrapopts: "",
					Architecture:    arch,
					Packages:        []rpm.Package{},
				})
			}
			cfg.DEB = debs

			cfg.DMG = config.DMG{
				Path:     "dmg",
				Packages: []string{},
			}

			flatpaks := []config.Flatpak{}
			for _, arch := range strings.Split(flatpakArchitectures, ",") {
				flatpaks = append(flatpaks, config.Flatpak{
					Path:         path.Join("flatpak", utils.GetArchIdentifier(arch)),
					Architecture: arch,
				})
			}
			cfg.Flatpak = flatpaks

			msis := []config.MSI{}
			for _, arch := range strings.Split(msiArchitectures, ",") {
				msis = append(msis, config.MSI{
					Path:         path.Join("msi", utils.GetArchIdentifier(arch)),
					Architecture: arch,
					Include:      `^\\b$`,
					Packages:     []string{},
				})
			}
			cfg.MSI = msis

			rpms := []config.RPM{}
			for _, arch := range strings.Split(rpmArchitectures, ",") {
				rpms = append(rpms, config.RPM{
					Path:         path.Join("rpm", "fedora", "40", utils.GetArchIdentifier(arch)),
					Trailer:      "fc40",
					Distro:       "fedora-40",
					Architecture: arch,
					Packages:     []rpm.Package{},
				})
			}
			cfg.RPM = rpms

			cfg.Binaries = config.Binaries{
				Path:     "binaries",
				Exclude:  binariesExclude,
				Packages: []string{},
			}
			cfg.Docs = config.Docs{
				Path: "docs",
			}

			b, err := yaml.Marshal(cfg)
			if err != nil {
				return err
			}

			if err := os.MkdirAll(dir, 0755); err != nil {
				return err
			}

			if err := os.WriteFile(filepath.Join(dir, "hydrapp.yaml"), b, 0664); err != nil {
				return err
			}
		}

		if err := os.WriteFile(filepath.Join(dir, "icon.png"), generators.IconTpl, 0664); err != nil {
			return err
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, "go.mod"),
			generators.GoModTpl,
			generators.GoModData{
				GoMod: goMod,
			},
		); err != nil {
			return err
		}

		switch projectTypeItems[projectTypeIndex].Name {
		case vanillaJSRESTKey:
			if err := generators.RenderTemplate(
				filepath.Join(dir, "main.go"),
				generators.GoMainVanillaJSRESTTpl,
				generators.GoMainData{
					GoMod: goMod,
				},
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, "android.go"),
				generators.AndroidVanillaJSRESTTpl,
				generators.AndroidData{
					GoMod:     goMod,
					JNIExport: strings.Replace(appID, ".", "_", -1),
				},
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, ".gitignore"),
				generators.GitignoreVanillaJSRESTTpl,
				nil,
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, "pkg", "backend", "server.go"),
				generators.BackendVanillaJSRESTTpl,
				nil,
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, "pkg", "frontend", "server.go"),
				generators.FrontendVanillaJSRESTTpl,
				nil,
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, "pkg", "frontend", "index.html"),
				generators.IndexHTMLVanillaJSRESTTpl,
				generators.IndexHTMLData{
					AppName: appName,
				},
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, "README.md"),
				generators.ReadmeMDVanillaJSRESTTpl,
				generators.ReadmeMDData{
					AppName:        appName,
					AppSummary:     appSummary,
					AppGitWeb:      strings.TrimSuffix(appGit, ".git"),
					AppDescription: appDescription,
					AppBaseURL:     appBaseurl,
					AppGit:         appGit,
					CurrentYear:    time.Now().Format("2006"),
					ReleaseAuthor:  releaseAuthor,
					LicenseSPDX:    licenseSPDX,
					Dir:            dir,
				},
			); err != nil {
				return err
			}
		case vanillaJSFormsKey:
			if err := generators.RenderTemplate(
				filepath.Join(dir, "main.go"),
				generators.GoMainVanillaJSFormsTpl,
				generators.GoMainData{
					GoMod: goMod,
				},
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, "android.go"),
				generators.AndroidVanillaJSFormsTpl,
				generators.AndroidData{
					GoMod:     goMod,
					JNIExport: strings.Replace(appID, ".", "_", -1),
				},
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, "pkg", "frontend", "server.go"),
				generators.FrontendVanillaJSFormsTpl,
				nil,
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, "pkg", "frontend", "index.html"),
				generators.IndexHTMLVanillaJSFormsTpl,
				generators.IndexHTMLData{
					AppName: appName,
				},
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, "README.md"),
				generators.ReadmeMDVanillaJSRESTTpl,
				generators.ReadmeMDData{
					AppName:        appName,
					AppSummary:     appSummary,
					AppGitWeb:      strings.TrimSuffix(appGit, ".git"),
					AppDescription: appDescription,
					AppBaseURL:     appBaseurl,
					AppGit:         appGit,
					CurrentYear:    time.Now().Format("2006"),
					ReleaseAuthor:  releaseAuthor,
					LicenseSPDX:    licenseSPDX,
					Dir:            dir,
				},
			); err != nil {
				return err
			}
		case reactPanrpcKey:
			if err := generators.RenderTemplate(
				filepath.Join(dir, "main.go"),
				generators.GoMainReactPanrpcTpl,
				generators.GoMainData{
					GoMod: goMod,
				},
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, "android.go"),
				generators.AndroidReactPanrpcTpl,
				generators.AndroidData{
					GoMod:     goMod,
					JNIExport: strings.Replace(appID, ".", "_", -1),
				},
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, ".gitignore"),
				generators.GitignoreReactPanrpcTpl,
				nil,
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, "pkg", "backend", "server.go"),
				generators.BackendReactPanrpcTpl,
				nil,
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, "pkg", "frontend", "server.go"),
				generators.FrontendReactPanrpcTpl,
				nil,
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, "pkg", "frontend", "src", "App.tsx"),
				generators.AppTSXTpl,
				generators.AppTSXData{
					AppName: appName,
				},
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, "pkg", "frontend", "src", "main.tsx"),
				generators.MainTSXTpl,
				nil,
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, "pkg", "frontend", "index.html"),
				generators.IndexHTMLReactPanrpcTpl,
				generators.IndexHTMLData{
					AppName: appName,
				},
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, "pkg", "frontend", "package.json"),
				generators.PackageJSONReactPanrpcTpl,
				generators.PackageJSONData{
					AppID:          appID,
					AppDescription: appDescription,
					ReleaseAuthor:  releaseAuthor,
					ReleaseEmail:   releaseEmail,
					LicenseSPDX:    licenseSPDX,
				},
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, "pkg", "frontend", "tsconfig.json"),
				generators.TsconfigJSONTpl,
				nil,
			); err != nil {
				return err
			}

			if err := generators.RenderTemplate(
				filepath.Join(dir, "README.md"),
				generators.ReadmeMDReactPanrpcTpl,
				generators.ReadmeMDData{
					AppName:        appName,
					AppSummary:     appSummary,
					AppGitWeb:      strings.TrimSuffix(appGit, ".git"),
					AppDescription: appDescription,
					AppBaseURL:     appBaseurl,
					AppGit:         appGit,
					CurrentYear:    time.Now().Format("2006"),
					ReleaseAuthor:  releaseAuthor,
					LicenseSPDX:    licenseSPDX,
					Dir:            dir,
				},
			); err != nil {
				return err
			}
		default:
			panic(errUnknownProjectType)
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, "LICENSE"),
			licenseText,
			nil,
		); err != nil {
			return err
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, "CODE_OF_CONDUCT.md"),
			generators.CodeOfConductMDTpl,
			generators.CodeOfConductMDData{
				ReleaseEmail: releaseEmail,
			},
		); err != nil {
			return err
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, ".github", "workflows", "hydrapp.yaml"),
			generators.HydrappYAMLTpl,
			generators.HydrappYAMLData{
				AppID:                         appID,
				ExperimentalGithubPagesAction: viper.GetBool(experimentalGithubPagesActionFlag),
			},
		); err != nil {
			return err
		}

		if !viper.GetBool(noNetworkFlag) {
			{
				log.Println("Downloading Go dependencies ...")

				cmd := exec.Command("go", "get", "-x", "./...")
				cmd.Dir = dir
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					return err
				}

				log.Println("Success!")
			}

			{
				log.Println("Generating Go dependencies ...")

				cmd := exec.Command("go", "generate", "-x", "./...")
				cmd.Dir = dir
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					return err
				}

				log.Println("Success!")
			}
		}

		{
			fmt.Printf(`Succesfully generated application. To start it, run the following:

cd %v
go run %v

You can find more information in the generated README.
`, dir, goMain)
		}

		return nil
	},
}

func init() {
	newCmd.PersistentFlags().Bool(noNetworkFlag, false, "Disable all network interaction")
	newCmd.PersistentFlags().Bool(experimentalGithubPagesActionFlag, false, "(Experimental) Use the GitHub Actions-based deploy strategy for GitHub pages instead of pushing to the gh-pages branch in the generated CI/CD configuration (disables support for publishing more than one hydrapp branch)")

	viper.AutomaticEnv()

	rootCmd.AddCommand(newCmd)
}
