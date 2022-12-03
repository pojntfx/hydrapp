package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/config"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/rpm"
	"github.com/pojntfx/hydrapp/hydrapp-generator/pkg/generators"
	"github.com/pojntfx/hydrapp/hydrapp-utils/pkg/utils"
	"gopkg.in/yaml.v3"
)

const (
	restKey      = "rest"
	formsKey     = "forms"
	dudirektaKey = "dudirekta"
)

var (
	errUnknownProjectType = errors.New("unknown project type")
)

func main() {
	noNetwork := flag.Bool("no-network", false, "Disable all network interaction")

	flag.Parse()

	projectTypeItems := []generators.ProjectTypeOption{
		{
			Name:        restKey,
			Description: "Simple starter project with a REST API to connect the frontend and backend",
		},
		{
			Name:        formsKey,
			Description: "Traditional starter project with Web 1.0-style forms to connect the frontend and backend",
		},
		{
			Name:        dudirektaKey,
			Description: "Complete starter project with bi-directional Dudirekta RPCs to connect the frontend and backend",
		},
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
		panic(err)
	}

	appID, err := (&promptui.Prompt{
		Label:   "App ID in reverse domain notation",
		Default: "com.github.example.myapp",
	}).Run()
	if err != nil {
		panic(err)
	}

	appName, err := (&promptui.Prompt{
		Label:   "App name",
		Default: "My App",
	}).Run()
	if err != nil {
		panic(err)
	}

	appSummary, err := (&promptui.Prompt{
		Label:   "App summary",
		Default: "My first app",
	}).Run()
	if err != nil {
		panic(err)
	}

	appDescription, err := (&promptui.Prompt{
		Label:   "App description",
		Default: "My first application, built with hydrapp.",
	}).Run()
	if err != nil {
		panic(err)
	}

	appHomepage, err := (&promptui.Prompt{
		Label:   "App homepage",
		Default: "https://github.com/example/myapp",
	}).Run()
	if err != nil {
		panic(err)
	}

	appGit, err := (&promptui.Prompt{
		Label:   "App git repo",
		Default: appHomepage + ".git",
	}).Run()
	if err != nil {
		panic(err)
	}

	appBaseurl, err := (&promptui.Prompt{
		Label:   "App base URL to expect the built assets to be published to",
		Default: "https://example.github.io/myapp/myapp/",
	}).Run()
	if err != nil {
		panic(err)
	}

	goMod, err := (&promptui.Prompt{
		Label:   "Go module name",
		Default: "github.com/example/myapp",
	}).Run()
	if err != nil {
		panic(err)
	}

	licenseSPDX, err := (&promptui.Prompt{
		Label:   "License SPDX identifier (see https://spdx.org/licenses/)",
		Default: "AGPL-3.0-or-later",
	}).Run()
	if err != nil {
		panic(err)
	}

	releaseAuthor, err := (&promptui.Prompt{
		Label:   "Release author name",
		Default: "Jean Doe",
	}).Run()
	if err != nil {
		panic(err)
	}

	releaseEmail, err := (&promptui.Prompt{
		Label:   "Release author email",
		Default: "jean.doe@example.com",
	}).Run()
	if err != nil {
		panic(err)
	}

	dir, err := (&promptui.Prompt{
		Label:   "Directory to write the app to",
		Default: "myapp",
	}).Run()
	if err != nil {
		panic(err)
	}

	_, advancedConfiguration, err := (&promptui.Select{
		Label: "Do you want to do any advanced configuration?",
		Items: []string{"no", "yes"},
	}).Run()
	if err != nil {
		panic(err)
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

	binariesExclude := "(android/*|ios/*|plan9/*|aix/*|linux/loong64|js/wasm)"

	if advancedConfiguration == "yes" {
		goMain, err = (&promptui.Prompt{
			Label:   "Go main package path",
			Default: goMain,
		}).Run()
		if err != nil {
			panic(err)
		}

		goFlags, err = (&promptui.Prompt{
			Label:   "Go flags to pass to the compiler",
			Default: goFlags,
		}).Run()
		if err != nil {
			panic(err)
		}

		goGenerate, err = (&promptui.Prompt{
			Label:   "Go generate command to run",
			Default: goGenerate,
		}).Run()
		if err != nil {
			panic(err)
		}

		goTests, err = (&promptui.Prompt{
			Label:   "Go test command to run",
			Default: goTests,
		}).Run()
		if err != nil {
			panic(err)
		}

		goImg, err = (&promptui.Prompt{
			Label:   "Go test OCI image to use",
			Default: goImg,
		}).Run()
		if err != nil {
			panic(err)
		}

		debArchitectures, err = (&promptui.Prompt{
			Label:   "DEB architectures to build for (comma-seperated list of GOARCH values)",
			Default: debArchitectures,
		}).Run()
		if err != nil {
			panic(err)
		}

		flatpakArchitectures, err = (&promptui.Prompt{
			Label:   "Flatpak architectures to build for (comma-seperated list of GOARCH values)",
			Default: flatpakArchitectures,
		}).Run()
		if err != nil {
			panic(err)
		}

		msiArchitectures, err = (&promptui.Prompt{
			Label:   "MSI architectures to build for (comma-seperated list of GOARCH values)",
			Default: msiArchitectures,
		}).Run()
		if err != nil {
			panic(err)
		}

		rpmArchitectures, err = (&promptui.Prompt{
			Label:   "RPM architectures to build for (comma-seperated list of GOARCH values)",
			Default: rpmArchitectures,
		}).Run()
		if err != nil {
			panic(err)
		}

		binariesExclude, err = (&promptui.Prompt{
			Label:   "Regex of binaries to exclude from compilation",
			Default: binariesExclude,
		}).Run()
		if err != nil {
			panic(err)
		}
	}

	licenseText := ""
	if !*noNetwork {
		log.Println("Fetching full license text from SPDX ...")

		res, err := http.Get("https://raw.githubusercontent.com/spdx/license-list-data/main/text/" + licenseSPDX + ".txt")
		if err != nil {
			panic(err)
		}
		if res.Body != nil {
			defer res.Body.Close()
		}
		if res.StatusCode != http.StatusOK {
			panic(errors.New(res.Status))
		}

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
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
				Path:         path.Join("rpm", "fedora", "37", utils.GetArchIdentifier(arch)),
				Trailer:      "1.fc37",
				Distro:       "fedora-37",
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
			panic(err)
		}

		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			panic(err)
		}

		if err := ioutil.WriteFile(filepath.Join(dir, "hydrapp.yaml"), b, os.ModePerm); err != nil {
			panic(err)
		}
	}

	if err := ioutil.WriteFile(filepath.Join(dir, "icon.png"), generators.IconTpl, os.ModePerm); err != nil {
		panic(err)
	}

	if err := generators.RenderTemplate(
		filepath.Join(dir, "go.mod"),
		generators.GoModTpl,
		generators.GoModData{
			GoMod: goMod,
		},
	); err != nil {
		panic(err)
	}

	switch projectTypeItems[projectTypeIndex].Name {
	case restKey:
		if err := generators.RenderTemplate(
			filepath.Join(dir, "main.go"),
			generators.GoMainRESTTpl,
			generators.GoMainData{
				GoMod: goMod,
			},
		); err != nil {
			panic(err)
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, "android.go"),
			generators.AndroidRESTTpl,
			generators.AndroidData{
				GoMod:     goMod,
				JNIExport: strings.Replace(appID, ".", "_", -1),
			},
		); err != nil {
			panic(err)
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, ".gitignore"),
			generators.GitignoreRESTTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, "pkg", "backend", "server.go"),
			generators.BackendRESTTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, "pkg", "frontend", "server.go"),
			generators.FrontendRESTTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, "pkg", "frontend", "index.html"),
			generators.IndexHTMLRESTTpl,
			generators.IndexHTMLData{
				AppName: appName,
			},
		); err != nil {
			panic(err)
		}
	case formsKey:
		if err := generators.RenderTemplate(
			filepath.Join(dir, "main.go"),
			generators.GoMainFormsTpl,
			generators.GoMainData{
				GoMod: goMod,
			},
		); err != nil {
			panic(err)
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, "android.go"),
			generators.AndroidFormsTpl,
			generators.AndroidData{
				GoMod:     goMod,
				JNIExport: strings.Replace(appID, ".", "_", -1),
			},
		); err != nil {
			panic(err)
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, "pkg", "frontend", "server.go"),
			generators.FrontendFormsTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, "pkg", "frontend", "index.html"),
			generators.IndexHTMLFormsTpl,
			generators.IndexHTMLData{
				AppName: appName,
			},
		); err != nil {
			panic(err)
		}
	case dudirektaKey:
		if err := generators.RenderTemplate(
			filepath.Join(dir, "main.go"),
			generators.GoMainDudirektaTpl,
			generators.GoMainData{
				GoMod: goMod,
			},
		); err != nil {
			panic(err)
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, "android.go"),
			generators.AndroidDudirektaTpl,
			generators.AndroidData{
				GoMod:     goMod,
				JNIExport: strings.Replace(appID, ".", "_", -1),
			},
		); err != nil {
			panic(err)
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, ".gitignore"),
			generators.GitignoreDudirektaTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, "pkg", "backend", "server.go"),
			generators.BackendDudirektaTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, "pkg", "frontend", "server.go"),
			generators.FrontendDudirektaTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, "pkg", "frontend", "src", "App.tsx"),
			generators.AppTSXTpl,
			generators.AppTSXData{
				AppName: appName,
			},
		); err != nil {
			panic(err)
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, "pkg", "frontend", "src", "main.tsx"),
			generators.MainTSXTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, "pkg", "frontend", "index.html"),
			generators.IndexHTMLDudirektaTpl,
			generators.IndexHTMLData{
				AppName: appName,
			},
		); err != nil {
			panic(err)
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, "pkg", "frontend", "package.json"),
			generators.PackageJSONTpl,
			generators.PackageJSONData{
				AppID:          appID,
				AppDescription: appDescription,
				ReleaseAuthor:  releaseAuthor,
				ReleaseEmail:   releaseEmail,
				LicenseSPDX:    licenseSPDX,
			},
		); err != nil {
			panic(err)
		}

		if err := generators.RenderTemplate(
			filepath.Join(dir, "pkg", "frontend", "tsconfig.json"),
			generators.TsconfigJSONTpl,
			nil,
		); err != nil {
			panic(err)
		}
	default:
		panic(errUnknownProjectType)
	}

	if err := generators.RenderTemplate(
		filepath.Join(dir, "LICENSE"),
		licenseText,
		nil,
	); err != nil {
		panic(err)
	}

	if err := generators.RenderTemplate(
		filepath.Join(dir, "CODE_OF_CONDUCT.md"),
		generators.CodeOfConductMDTpl,
		generators.CodeOfConductMDData{
			ReleaseEmail: releaseEmail,
		},
	); err != nil {
		panic(err)
	}

	if err := generators.RenderTemplate(
		filepath.Join(dir, "README.md"),
		generators.ReadmeMDTpl,
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
		},
	); err != nil {
		panic(err)
	}

	if err := generators.RenderTemplate(
		filepath.Join(dir, ".github", "workflows", "hydrapp.yaml"),
		generators.HydrappYAMLTpl,
		generators.HydrappYAMLData{
			AppID: appID,
		},
	); err != nil {
		panic(err)
	}

	if !*noNetwork {
		{
			log.Println("Downloading Go dependencies ...")

			cmd := exec.Command("go", "get", "-x", "./...")
			cmd.Dir = dir
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				panic(err)
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
				panic(err)
			}

			log.Println("Success!")
		}

		{
			fmt.Printf(`Succesfully generated application. To start it, run the following:

cd %v
go run %v

You can find more information in the generated README.
`, dir, goMain)
		}
	}
}
