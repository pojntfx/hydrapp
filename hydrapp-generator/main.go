package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	_ "embed"

	"github.com/manifoldco/promptui"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/config"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/rpm"
	"github.com/pojntfx/hydrapp/hydrapp-utils/pkg/utils"
	"gopkg.in/yaml.v3"
)

const (
	restKey      = "rest"
	formsKey     = "forms"
	dudirektaKey = "dudirekta"
)

var (
	//go:embed icon.png.tpl
	iconTpl []byte

	//go:embed go.mod.tpl
	goModTpl string

	//go:embed main_dudirekta.go.tpl
	goMainDudirektaTpl string

	//go:embed main_forms.go.tpl
	goMainFormsTpl string

	//go:embed main_rest.go.tpl
	goMainRESTTpl string

	//go:embed android_dudirekta.go.tpl
	androidDudirektaTpl string

	//go:embed android_forms.go.tpl
	androidFormsTpl string

	//go:embed android_rest.go.tpl
	androidRESTTpl string

	//go:embed .gitignore_dudirekta.tpl
	gitignoreDudirektaTpl string

	//go:embed .gitignore_rest.tpl
	gitignoreRESTTpl string

	//go:embed backend_dudirekta.go.tpl
	backendDudirektaTpl string

	//go:embed backend_rest.go.tpl
	backendRESTTpl string

	//go:embed frontend_dudirekta.go.tpl
	frontendDudirektaTpl string

	//go:embed frontend_rest.go.tpl
	frontendRESTTpl string

	//go:embed frontend_forms.go.tpl
	frontendFormsTpl string

	//go:embed App.tsx.tpl
	appTSXTpl string

	//go:embed main.tsx.tpl
	mainTSXTpl string

	//go:embed index_dudirekta.html.tpl
	indexHTMLDudirektaTpl string

	//go:embed index_rest.html.tpl
	indexHTMLRESTTpl string

	//go:embed index_forms.html.tpl
	indexHTMLFormsTpl string

	//go:embed package.json.tpl
	packageJSONTpl string

	//go:embed tsconfig.json.tpl
	tsconfigJSONTpl string

	errUnknownProjectType = errors.New("unknown project type")
)

type goModData struct {
	GoMod string
}

type goMainData struct {
	GoMod string
}

type androidData struct {
	GoMod     string
	JNIExport string
}

type appTSXData struct {
	AppName string
}

type indexHTMLData struct {
	AppName string
}

type packageJSONData struct {
	AppID          string
	AppDescription string
	ReleaseAuthor  string
	ReleaseEmail   string
	LicenseSPDX    string
}

func renderTemplate(path string, tpl string, data any) error {
	// Assume that templates without data are just files
	if data == nil {
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}

		return ioutil.WriteFile(path, []byte(tpl), os.ModePerm)
	}

	t, err := template.New(path).Parse(tpl)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(t.Name()), os.ModePerm); err != nil {
		return err
	}

	dst, err := os.OpenFile(t.Name(), os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer dst.Close()

	return t.Execute(dst, data)
}

func main() {
	goMain := flag.String("go-main", ".", "Go main package path")
	goFlags := flag.String("go-flags", "", "Go flags to pass to the compiler")
	goGenerate := flag.String("go-generate", "go generate ./...", "Go generate command to run")
	goTests := flag.String("go-tests", "go test ./...", "Go test command to run")
	goImg := flag.String("go-img", "ghcr.io/pojntfx/hydrapp-build-tests:main", "Go test OCI image to use")

	debArchitectures := flag.String("deb-architectures", "amd64", "DEB architectures to build for (comma-seperated list of GOARCH values)")
	flatpakArchitectures := flag.String("flatpak-architectures", "amd64", "Flatpak architectures to build for (comma-seperated list of GOARCH values)")
	msiArchitectures := flag.String("msi-architectures", "amd64", "MSI architectures to build for (comma-seperated list of GOARCH values)")
	rpmArchitectures := flag.String("rpm-architectures", "amd64", "RPM architectures to build for (comma-seperated list of GOARCH values)")

	binariesExclude := flag.String("binaries-exclude", "(android/*|ios/*|plan9/*|aix/*|linux/loong64|js/wasm)", "Regex of binaries to exclude from compilation")

	noNetwork := flag.Bool("no-network", false, "Disable all network interaction")

	flag.Parse()

	// TODO: Add help menu for each select item
	_, projectType, err := (&promptui.Select{
		Label: "Project type to generate",
		Items: []string{restKey, formsKey, dudirektaKey},
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

	dir, err := (&promptui.Prompt{
		Label:   "Directory to write the app to",
		Default: "myapp",
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

	_, advancedConfiguration, err := (&promptui.Select{
		Label: "Do you want to do any advanced configuration?",
		Items: []string{"no", "yes"},
	}).Run()
	if err != nil {
		panic(err)
	}

	if advancedConfiguration == "yes" {

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
			Main:     *goMain,
			Flags:    *goFlags,
			Generate: *goGenerate,
			Tests:    *goTests,
			Image:    *goImg,
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
		for _, arch := range strings.Split(*debArchitectures, ",") {
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
		for _, arch := range strings.Split(*flatpakArchitectures, ",") {
			flatpaks = append(flatpaks, config.Flatpak{
				Path:         path.Join("flatpak", utils.GetArchIdentifier(arch)),
				Architecture: arch,
			})
		}
		cfg.Flatpak = flatpaks

		msis := []config.MSI{}
		for _, arch := range strings.Split(*msiArchitectures, ",") {
			msis = append(msis, config.MSI{
				Path:         path.Join("msi", utils.GetArchIdentifier(arch)),
				Architecture: arch,
				Include:      `^\\b$`,
				Packages:     []string{},
			})
		}
		cfg.MSI = msis

		rpms := []config.RPM{}
		for _, arch := range strings.Split(*rpmArchitectures, ",") {
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
			Exclude:  *binariesExclude,
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

	if err := ioutil.WriteFile(filepath.Join(dir, "icon.png"), iconTpl, os.ModePerm); err != nil {
		panic(err)
	}

	if err := renderTemplate(
		filepath.Join(dir, "go.mod"),
		goModTpl,
		goModData{
			GoMod: goMod,
		},
	); err != nil {
		panic(err)
	}

	switch projectType {
	case restKey:
		if err := renderTemplate(
			filepath.Join(dir, "main.go"),
			goMainRESTTpl,
			goMainData{
				GoMod: goMod,
			},
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(dir, "android.go"),
			androidRESTTpl,
			androidData{
				GoMod:     goMod,
				JNIExport: strings.Replace(appID, ".", "_", -1),
			},
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(dir, ".gitignore"),
			gitignoreRESTTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(dir, "pkg", "backend", "server.go"),
			backendRESTTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(dir, "pkg", "frontend", "server.go"),
			frontendRESTTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(dir, "pkg", "frontend", "index.html"),
			indexHTMLRESTTpl,
			indexHTMLData{
				AppName: appName,
			},
		); err != nil {
			panic(err)
		}
	case formsKey:
		if err := renderTemplate(
			filepath.Join(dir, "main.go"),
			goMainFormsTpl,
			goMainData{
				GoMod: goMod,
			},
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(dir, "android.go"),
			androidFormsTpl,
			androidData{
				GoMod:     goMod,
				JNIExport: strings.Replace(appID, ".", "_", -1),
			},
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(dir, "pkg", "frontend", "server.go"),
			frontendFormsTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(dir, "pkg", "frontend", "index.html"),
			indexHTMLFormsTpl,
			indexHTMLData{
				AppName: appName,
			},
		); err != nil {
			panic(err)
		}
	case dudirektaKey:
		if err := renderTemplate(
			filepath.Join(dir, "main.go"),
			goMainDudirektaTpl,
			goMainData{
				GoMod: goMod,
			},
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(dir, "android.go"),
			androidDudirektaTpl,
			androidData{
				GoMod:     goMod,
				JNIExport: strings.Replace(appID, ".", "_", -1),
			},
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(dir, ".gitignore"),
			gitignoreDudirektaTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(dir, "pkg", "backend", "server.go"),
			backendDudirektaTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(dir, "pkg", "frontend", "server.go"),
			frontendDudirektaTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(dir, "pkg", "frontend", "src", "App.tsx"),
			appTSXTpl,
			appTSXData{
				AppName: appName,
			},
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(dir, "pkg", "frontend", "src", "main.tsx"),
			mainTSXTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(dir, "pkg", "frontend", "index.html"),
			indexHTMLDudirektaTpl,
			indexHTMLData{
				AppName: appName,
			},
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(dir, "pkg", "frontend", "package.json"),
			packageJSONTpl,
			packageJSONData{
				AppID:          appID,
				AppDescription: appDescription,
				ReleaseAuthor:  releaseAuthor,
				ReleaseEmail:   releaseEmail,
				LicenseSPDX:    licenseSPDX,
			},
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(dir, "pkg", "frontend", "tsconfig.json"),
			tsconfigJSONTpl,
			nil,
		); err != nil {
			panic(err)
		}
	default:
		panic(errUnknownProjectType)
	}

	if err := renderTemplate(
		filepath.Join(dir, "LICENSE"),
		licenseText,
		nil,
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
	}
}
