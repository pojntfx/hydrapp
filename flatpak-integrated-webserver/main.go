//go:build !android
// +build !android

package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pojntfx/multi-browser-electron/flatpak-integrated-webserver/pkg/backend"
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

const (
	name = "Integrated Webserver Example Flatpak"
	id   = "com.pojtinger.flatpakIntegratedWebserver"

	spawnCmd  = "flatpak-spawn"
	spawnHost = "--host"
	whichCmd  = "which"
)

func main() {
	// Start the integrated webserver server
	url, stop, err := backend.StartServer()
	if err != nil {
		log.Fatalln("could not start integrated webserver:", err)
	}
	defer stop()

	// Use the user-prefered browser if specified
	browser := []string{os.Getenv("HYDRAPP_BROWSER")}

	// Check if we are in flatpak
	runningInFlatpak := false
	if _, err := exec.LookPath(spawnCmd); err == nil {
		runningInFlatpak = true
	}

	// Find supported browser
	if browser[0] == "" {
		for _, knownBrowser := range knownBrowsers {
			if runningInFlatpak {
				// Find supported browser in Flatpak install
				if err := exec.Command(spawnCmd, spawnHost, whichCmd, knownBrowser).Run(); err == nil {
					browser = []string{spawnCmd, spawnHost, knownBrowser}

					break
				}
			} else {
				// Find supported browser in native install
				if _, err := exec.LookPath(knownBrowser); err == nil {
					browser = []string{knownBrowser}

					break
				}
			}
		}
	} else {
		if runningInFlatpak {
			// Allow same override in Flatpak by spawning the prefered browser on the host
			browser = append([]string{spawnCmd, spawnHost}, browser...)
		}
	}

	if browser[0] == "" {
		log.Fatalf("could not find a supported browser, tried preferred browser (set with the HYDRAPP_BROWSER env variable) \"%v\" and known browsers \"%v\"", browser, knownBrowsers)
	}

	// Create a profile for the app
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("could not get user's config dir:", err)
	}
	userDataDir := filepath.Join(userConfigDir, id)

	// Create the browser instance
	execLine := append(
		browser,
		append(
			[]string{
				"--app=" + url,
				"--class=" + name,
				"--user-data-dir=" + userDataDir,
				"--no-first-run",
				"--no-default-browser-check",
			},
			os.Args[1:]...,
		)...,
	)
	cmd := exec.Command(
		execLine[0],
		execLine[1:]...,
	)

	// Use systemd stdout, stderr and stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Start the browser and wait for the user to close it
	if err := cmd.Run(); err != nil {
		log.Fatal("could not launch browser:", err)
	}
}
