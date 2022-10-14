package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/apk"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/deb"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/dmg"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/flatpak"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/msi"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/rpm"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/xdg"
)

func main() {
	appID := flag.String("app-id", "com.pojtinger.felicitas.hydrapp.example", "App ID")
	appName := flag.String("app-name", "Hydrapp Example", "App name")
	appSummary := flag.String("app-summary", "Hydrapp example app", "App summary")
	appDescription := flag.String("app-description", "A simple Hydrapp example app.", "App description")
	appURL := flag.String("app-url", "https://github.com/pojntfx/hydrapp/tree/main/hydrapp-example", "App URL")
	appSPDX := flag.String("app-spdx", "AGPL-3.0", "App SPDX license identifier")
	appReleases := flag.String("app-releases", `[{ "version": "0.0.1", "date": "2022-10-11", "description": "Initial release", "author": "Felicitas Pojtinger", "email": "felicitas@pojtinger.com" }]`, "App SPDX license identifier")
	extraRHELPackages := flag.String("extra-rhel-packages", `[]`, `Extra RHEL packages (in format { "name": "firefox", "version": "89" })`)
	extraSUSEPackages := flag.String("extra-suse-packages", `[]`, `Extra SUSE packages (in format { "name": "firefox", "version": "89" })`)
	appBackendPkg := flag.String("app-backend-pkg", "github.com/pojntfx/hydrapp/hydrapp-example/pkg/backend", "App backend package")
	appFrontendPkg := flag.String("app-frontend-pkg", "github.com/pojntfx/hydrapp/hydrapp-example/pkg/frontend", "App frontend package")

	flag.Parse()

	var releases []rpm.Release
	if err := json.Unmarshal([]byte(*appReleases), &releases); err != nil {
		panic(err)
	}

	var rhelPackages []rpm.Package
	if err := json.Unmarshal([]byte(*extraRHELPackages), &rhelPackages); err != nil {
		panic(err)
	}

	var susePackages []rpm.Package
	if err := json.Unmarshal([]byte(*extraSUSEPackages), &susePackages); err != nil {
		panic(err)
	}

	for _, renderer := range []*renderers.Renderer{
		apk.NewManifestRenderer(
			*appID,
			*appName,
		),
		apk.NewActivityRenderer(
			*appID,
		),
		apk.NewHeaderRenderer(),
		apk.NewBindingsRenderer(
			*appID,
			*appBackendPkg,
			*appFrontendPkg,
		),
		apk.NewImplementationRenderer(),
		xdg.NewDesktopRenderer(
			*appID,
			*appName,
			*appDescription,
		),
		xdg.NewMetainfoRenderer(
			*appID,
			*appName,
			*appDescription,
			*appSummary,
			*appSPDX,
			*appURL,
			releases,
		),
		rpm.NewSpecRenderer(
			*appID,
			*appName,
			*appDescription,
			*appSummary,
			*appSPDX,
			*appURL,
			releases,
			rhelPackages,
			susePackages,
		),
		msi.NewWixRenderer(
			*appID,
			*appName,
			*appDescription,
			releases,
		),
		flatpak.NewManifestRenderer(
			*appID,
		),
		flatpak.NewSdkRenderer(),
		dmg.NewInfoRenderer(
			*appID,
			*appName,
			releases,
		),
		deb.NewChangelogRenderer(
			*appID,
			releases,
		),
		deb.NewCompatRenderer(),
		deb.NewFormatRenderer(),
		deb.NewManpagesRenderer(
			*appID,
		),
		deb.NewOptionsRenderer(),
	} {
		if path, content, err := renderer.Render(); err != nil {
			panic(err)
		} else {
			fmt.Printf("%v\n%v\n", path, content)
		}
	}
}
