//go:build !android

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/pojntfx/multi-browser-electron/desktop-integrated-webserver/pkg/backend"
)

var knownBrowsers = []string{
	"google-chrome",
	"google-chrome-stable",
	"google-chrome-beta",
	"google-chrome-unstable",
	"brave-browser",
	"brave-browser-stable",
	"brave-browser-beta",
	"brave-browser-nightly",
	"microsoft-edge",
	"microsoft-edge-beta",
	"microsoft-edge-dev",
	"ungoogled-chromium",
	"chromium-browser",
}

func main() {
	url, done, stop, err := backend.StartServer()
	if err != nil {
		log.Fatalln("could not start integrated webserver:", err)
	}

	name := "Integrated Webserver Example"
	// id := "com.pojtinger.felicitas.integratedWebserverExample"

	browser := os.Getenv("HYDRAPP_BROWSER")
	if browser == "" {
		for _, knownBrowser := range knownBrowsers {
			if _, err := exec.LookPath(knownBrowser); err == nil {
				browser = knownBrowser

				break
			}
		}

		if browser == "" {
			log.Fatalf("could not find a supported browser, tried preferred browser (set with the HYDRAPP_BROWSER env variable) \"%v\" and known browsers \"%v\"", browser, knownBrowsers)
		}
	}

	// This will be used for non-Chromium based browsers
	// userConfigDir, err := os.UserConfigDir()
	// if err != nil {
	// 	log.Fatal("could not get user's config dir:", err)
	// }
	// userDataDir := filepath.Join(userConfigDir, id)

	output, err := exec.Command(
		browser,
		append(
			[]string{
				"--app=" + url,
				"--class=" + name,
				"--no-first-run",
				"--no-default-browser-check",
			},
			os.Args[1:]...,
		)...,
	).CombinedOutput()
	if err != nil {
		stop()

		log.Fatal("could not launch browser:", string(output)+":", err)
	}

	fmt.Print(string(output))

	<-done
}
