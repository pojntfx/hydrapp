//go:build !android
// +build !android

package main

import (
	"context"
	"os"

	"github.com/pojntfx/hydrapp/example/pkg/backend"
	"github.com/pojntfx/hydrapp/example/pkg/browser"
	_ "github.com/pojntfx/hydrapp/example/pkg/fixes"
	"github.com/pojntfx/hydrapp/example/pkg/update"
	"github.com/pojntfx/hydrapp/example/pkg/utils"
)

const (
	appName    = "Hydrapp Example"                         // App name
	appID      = "com.pojtinger.felicitas.hydrapp.example" // App ID
	appVersion = "v0.0.1"                                  // App version

	updateAPIURL = "https://api.github.com/" // GitHub/Gitea API endpoint to use
	updateOwner  = "pojntfx"                 // Repository owner
	updateRepo   = "hydrapp"                 // Repository name
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Apply the self-update if not disabled
	browserState := &update.BrowserState{}
	if os.Getenv("HYDRAPP_SELFUPDATE") != "false" {
		go func() {
			if err := update.Update(
				ctx,

				updateAPIURL,
				updateOwner,
				updateRepo,

				appVersion,
				appID,

				browserState,
				func(msg string, err error) {
					utils.HandlePanic(appName, "could not check for updates (disable it by setting the HYDRAPP_SELFUPDATE env variable to false): "+msg, err)
				},
			); err != nil {
				utils.HandlePanic(appName, "could not check for updates (disable it by setting the HYDRAPP_SELFUPDATE env variable to false)", err)
			}
		}()
	}

	// Start the integrated webserver server
	url, stop, err := backend.StartServer()
	if err != nil {
		utils.HandlePanic(appName, "could not start integrated webserver", err)
	}
	defer stop()

	// Use the user-prefered browser binary and type if specified
	browserBinary := os.Getenv("HYDRAPP_BROWSER")
	browserType := os.Getenv("HYDRAPP_TYPE")

	browser.LaunchBrowser(
		url,
		appName,
		appID,

		browserBinary,
		browserType,

		browser.ChromiumLikeBrowsers,
		browser.FirefoxLikeBrowsers,
		browser.EpiphanyLikeBrowsers,

		browserState,
		func(msg string, err error) {
			utils.HandlePanic(appName, msg, err)
		},
	)
}
