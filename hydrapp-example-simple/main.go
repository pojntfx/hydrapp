//go:build !android
// +build !android

package main

import (
	"context"
	"log"
	"os"

	"github.com/pojntfx/hydrapp/hydrapp-example-simple/pkg/backend"
	"github.com/pojntfx/hydrapp/hydrapp-example-simple/pkg/frontend"
	"github.com/pojntfx/hydrapp/hydrapp-utils/pkg/browser"
	_ "github.com/pojntfx/hydrapp/hydrapp-utils/pkg/fixes"
	"github.com/pojntfx/hydrapp/hydrapp-utils/pkg/update"
	"github.com/pojntfx/hydrapp/hydrapp-utils/pkg/utils"
)

const (
	appName    = "Hydrapp Simple Example"                         // App name
	appID      = "com.pojtinger.felicitas.hydrapp.example.simple" // App ID
	appVersion = "v0.0.1"                                         // App version

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

	// Start the backend
	backendURL, stopBackend, err := backend.StartServer(os.Getenv("HYDRAPP_BACKEND_LADDR"), true)
	if err != nil {
		utils.HandlePanic(appName, "could not start backend", err)
	}
	defer stopBackend()

	log.Println("Backend URL:", backendURL)

	// Start the frontend
	frontendURL, stopFrontend, err := frontend.StartServer(os.Getenv("HYDRAPP_FRONTEND_LADDR"), backendURL, true)
	if err != nil {
		utils.HandlePanic(appName, "could not start frontend", err)
	}
	defer stopFrontend()

	log.Println("Frontend URL:", frontendURL)

	browser.LaunchBrowser(
		frontendURL,
		appName,
		appID,

		os.Getenv("HYDRAPP_BROWSER"),
		os.Getenv("HYDRAPP_TYPE"),

		browser.ChromiumLikeBrowsers,
		browser.FirefoxLikeBrowsers,
		browser.EpiphanyLikeBrowsers,

		browserState,
		func(msg string, err error) {
			utils.HandlePanic(appName, msg, err)
		},
	)
}
