package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/androidmanifest"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/desktopentry"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/flatpak"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/metainfo"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/spec"
	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/renderers/wix"
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

	flag.Parse()

	var releases []spec.Release
	if err := json.Unmarshal([]byte(*appReleases), &releases); err != nil {
		panic(err)
	}

	var rhelPackages []spec.Package
	if err := json.Unmarshal([]byte(*extraRHELPackages), &rhelPackages); err != nil {
		panic(err)
	}

	var susePackages []spec.Package
	if err := json.Unmarshal([]byte(*extraSUSEPackages), &susePackages); err != nil {
		panic(err)
	}

	for _, renderer := range []*renderers.Renderer{
		androidmanifest.NewRenderer(
			*appID,
			*appName,
		),
		desktopentry.NewRenderer(
			*appID,
			*appName,
			*appDescription,
		),
		metainfo.NewRenderer(
			*appID,
			*appName,
			*appDescription,
			*appSummary,
			*appSPDX,
			*appURL,
			releases,
		),
		spec.NewRenderer(
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
		wix.NewRenderer(
			*appID,
			*appName,
			*appDescription,
			releases,
		),
		flatpak.NewRenderer(
			*appID,
		),
	} {
		if path, content, err := renderer.Render(); err != nil {
			panic(err)
		} else {
			fmt.Printf("%v\n%v\n", path, content)
		}
	}
}
