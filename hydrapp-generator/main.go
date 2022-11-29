package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	_ "embed"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/config"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/rpm"
	"github.com/pojntfx/hydrapp/hydrapp-utils/pkg/utils"
	"gopkg.in/yaml.v3"
)

const (
	agplv3LicenseTextSummary = `This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.
.
You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
.
On Debian systems, the complete text of the GNU General
Public License version 3 can be found in "/usr/share/common-licenses/GPL-3".`
)

var (
	//go:embed icon.png.tpl
	iconTpl []byte

	//go:embed go.mod.tpl
	goModTpl string

	//go:embed main_full.go.tpl
	goMainFullTpl string

	//go:embed main_forms.go.tpl
	goMainFormsTpl string

	//go:embed android_full.go.tpl
	androidFullTpl string

	//go:embed android_forms.go.tpl
	androidFormsTpl string

	//go:embed .gitignore_full.tpl
	gitignoreFullTpl string

	//go:embed backend_full.go.tpl
	backendFullTpl string

	//go:embed frontend_full.go.tpl
	frontendFullTpl string

	//go:embed frontend_forms.go.tpl
	frontendFormsTpl string

	//go:embed App.tsx.tpl
	appTSXTpl string

	//go:embed main.tsx.tpl
	mainTSXTpl string

	//go:embed index_full.html.tpl
	indexHTMLFullTpl string

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
	appID := flag.String("app-id", "com.github.example.myapp", "App ID in reverse domain notation")
	appName := flag.String("app-name", "My App", "App name")
	appSummary := flag.String("app-summary", "My first app", "App summary")
	appDescription := flag.String("app-description", "My first application, built with hydrapp.", "App description")
	appHomepage := flag.String("app-homepage", "https://github.com/example/myapp", "App homepage")
	appGit := flag.String("app-git", "https://github.com/example/myapp.git", "App git repo")
	appBaseurl := flag.String("app-baseurl", "https://example.github.io/myapp/myapp/", "App base URL to expect the built assets to be published to")

	goMain := flag.String("go-main", ".", "Go main package path")
	goFlags := flag.String("go-flags", "", "Go flags to pass to the compiler")
	goGenerate := flag.String("go-generate", "go generate ./...", "Go generate command to run")
	goTests := flag.String("go-tests", "go test ./...", "Go test command to run")
	goImg := flag.String("go-img", "ghcr.io/pojntfx/hydrapp-build-tests:main", "Go test OCI image to use")
	goMod := flag.String("go-mod", "github.com/example/myapp", "Go module name")

	licenseSPDX := flag.String("license-spdx", "AGPL-3.0", "License SPDX identifier (see https://spdx.org/licenses/)")
	licenseTextSummary := flag.String("license-text-summary", agplv3LicenseTextSummary, "License summary text")

	releaseAuthor := flag.String("release-author", "Jean Doe", "Release author name")
	releaseEmail := flag.String("release-email", "jean.doe@example.com", "Release author email")

	debArchitectures := flag.String("deb-architectures", "amd64", "DEB architectures to build for (comma-seperated list of GOARCH values)")
	flatpakArchitectures := flag.String("flatpak-architectures", "amd64", "Flatpak architectures to build for (comma-seperated list of GOARCH values)")
	msiArchitectures := flag.String("msi-architectures", "amd64", "MSI architectures to build for (comma-seperated list of GOARCH values)")
	rpmArchitectures := flag.String("rpm-architectures", "amd64", "RPM architectures to build for (comma-seperated list of GOARCH values)")

	binariesExclude := flag.String("binaries-exclude", "(android/*|ios/*|plan9/*|aix/*|linux/loong64|js/wasm)", "Regex of binaries to exclude from compilation")

	dir := flag.String("dir", "myapp", "Directory to write the app to")

	projectType := flag.String("project-type", "full", "Project type to generate (simple, forms or full)")

	flag.Parse()

	{
		cfg := config.Root{}
		cfg.App = config.App{
			ID:          *appID,
			Name:        *appName,
			Summary:     *appSummary,
			Description: *appDescription,
			Homepage:    *appHomepage,
			Git:         *appGit,
			BaseURL:     *appBaseurl,
		}
		cfg.Go = config.Go{
			Main:     *goMain,
			Flags:    *goFlags,
			Generate: *goGenerate,
			Tests:    *goTests,
			Image:    *goImg,
		}
		cfg.License = config.License{
			SPDX: *licenseSPDX,
			Text: *licenseTextSummary,
		}
		cfg.Releases = []renderers.Release{
			{
				Version:     "0.0.1",
				Date:        time.Now(),
				Description: "Initial release",
				Author:      *releaseAuthor,
				Email:       *releaseEmail,
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

		if err := os.MkdirAll(*dir, os.ModePerm); err != nil {
			panic(err)
		}

		if err := ioutil.WriteFile(filepath.Join(*dir, "hydrapp.yaml"), b, os.ModePerm); err != nil {
			panic(err)
		}
	}

	if err := ioutil.WriteFile(filepath.Join(*dir, "icon.png"), iconTpl, os.ModePerm); err != nil {
		panic(err)
	}

	if err := renderTemplate(
		filepath.Join(*dir, "go.mod"),
		goModTpl,
		goModData{
			GoMod: *goMod,
		},
	); err != nil {
		panic(err)
	}

	switch *projectType {
	case "full":
		if err := renderTemplate(
			filepath.Join(*dir, "main.go"),
			goMainFullTpl,
			goMainData{
				GoMod: *goMod,
			},
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(*dir, "android.go"),
			androidFullTpl,
			androidData{
				GoMod:     *goMod,
				JNIExport: strings.Replace(*appID, ".", "_", -1),
			},
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(*dir, ".gitignore"),
			gitignoreFullTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(*dir, "pkg", "backend", "server.go"),
			backendFullTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(*dir, "pkg", "frontend", "server.go"),
			frontendFullTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(*dir, "pkg", "frontend", "src", "App.tsx"),
			appTSXTpl,
			appTSXData{
				AppName: *appName,
			},
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(*dir, "pkg", "frontend", "src", "main.tsx"),
			mainTSXTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(*dir, "pkg", "frontend", "index.html"),
			indexHTMLFullTpl,
			indexHTMLData{
				AppName: *appName,
			},
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(*dir, "pkg", "frontend", "package.json"),
			packageJSONTpl,
			packageJSONData{
				AppID:          *appID,
				AppDescription: *appDescription,
				ReleaseAuthor:  *releaseAuthor,
				ReleaseEmail:   *releaseEmail,
				LicenseSPDX:    *licenseSPDX,
			},
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(*dir, "pkg", "frontend", "tsconfig.json"),
			tsconfigJSONTpl,
			nil,
		); err != nil {
			panic(err)
		}
	case "forms":
		if err := renderTemplate(
			filepath.Join(*dir, "main.go"),
			goMainFormsTpl,
			goMainData{
				GoMod: *goMod,
			},
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(*dir, "android.go"),
			androidFormsTpl,
			androidData{
				GoMod:     *goMod,
				JNIExport: strings.Replace(*appID, ".", "_", -1),
			},
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(*dir, "pkg", "frontend", "server.go"),
			frontendFormsTpl,
			nil,
		); err != nil {
			panic(err)
		}

		if err := renderTemplate(
			filepath.Join(*dir, "pkg", "frontend", "index.html"),
			indexHTMLFormsTpl,
			indexHTMLData{
				AppName: *appName,
			},
		); err != nil {
			panic(err)
		}
	default:
		panic(errUnknownProjectType)
	}

	if err := renderTemplate(
		filepath.Join(*dir, "LICENSE"),
		*licenseTextSummary,
		nil,
	); err != nil {
		panic(err)
	}

	{
		cmd := exec.Command("go", "get", "-x", "./...")
		cmd.Dir = *dir
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}

	{
		cmd := exec.Command("go", "generate", "-x", "./...")
		cmd.Dir = *dir
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}
}
