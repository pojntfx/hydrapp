package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/config"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"gopkg.in/yaml.v3"
)

const (
	agplv3LicenseText = `This program is free software: you can redistribute it and/or modify
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

	licenseSPDX := flag.String("license-spdx", "AGPL-3.0", "License SPDX identifier (see https://spdx.org/licenses/)")
	licenseText := flag.String("license-text", agplv3LicenseText, "License summary text")

	releaseAuthor := flag.String("release-author", "Jean Doe", "Release author name")
	releaseEmail := flag.String("release-email", "jean.doe@example.com", "Release author email")

	// debArchitectures := flag.String("deb-architectures", "amd64", "DEB architectures to build for (comma-seperated list of GOARCH values)")
	// flatpakArchitectures := flag.String("flatpak-architectures", "amd64", "Flatpak architectures to build for (comma-seperated list of GOARCH values)")
	// msiArchitectures := flag.String("msi-architectures", "amd64", "MSI architectures to build for (comma-seperated list of GOARCH values)")
	// rpmArchitectures := flag.String("rpm-architectures", "amd64", "RPM architectures to build for (comma-seperated list of GOARCH values)")

	// binariesExclude := flag.String("binaries-exclude", "(android/*|ios/*|plan9/*|aix/*|linux/loong64|js/wasm)", "Regex of binaries to exclude from compilation")

	dir := flag.String("dir", "myapp", "Directory to write the app to")

	flag.Parse()

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
		Text: *licenseText,
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
