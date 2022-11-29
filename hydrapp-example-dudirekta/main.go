//go:build !android
// +build !android

package main

import (
	"context"
	_ "embed"
	"log"
	"os"
	"time"

	"github.com/pojntfx/hydrapp/hydrapp-builder/pkg/config"
	"github.com/pojntfx/hydrapp/hydrapp-example-dudirekta/pkg/backend"
	"github.com/pojntfx/hydrapp/hydrapp-example-dudirekta/pkg/frontend"
	"github.com/pojntfx/hydrapp/hydrapp-utils/pkg/browser"
	_ "github.com/pojntfx/hydrapp/hydrapp-utils/pkg/fixes"
	"github.com/pojntfx/hydrapp/hydrapp-utils/pkg/update"
	"github.com/pojntfx/hydrapp/hydrapp-utils/pkg/utils"
)

//go:embed hydrapp.yaml
var configFile []byte

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Parse(configFile)
	if err != nil {
		utils.HandlePanic("App", "could not parse config file", err)

		return
	}

	// Apply the self-update
	browserState := &update.BrowserState{}
	go update.Update(
		ctx,

		cfg,
		browserState,
		utils.HandlePanic,
	)

	// Start the backend
	backendURL, stopBackend, err := backend.StartServer(ctx, os.Getenv(utils.EnvBackendLaddr), time.Second*10, true)
	if err != nil {
		utils.HandlePanic(cfg.App.Name, "could not start backend", err)
	}
	defer stopBackend()

	log.Println("Backend URL:", backendURL)

	// Start the frontend
	frontendURL, stopFrontend, err := frontend.StartServer(ctx, os.Getenv(utils.EnvFrontendLaddr), backendURL, true)
	if err != nil {
		utils.HandlePanic(cfg.App.Name, "could not start frontend", err)
	}
	defer stopFrontend()

	log.Println("Frontend URL:", frontendURL)

	browser.LaunchBrowser(
		frontendURL,
		cfg.App.Name,
		cfg.App.ID,

		os.Getenv(utils.EnvBrowser),
		os.Getenv(utils.EnvType),

		browser.ChromiumLikeBrowsers,
		browser.FirefoxLikeBrowsers,
		browser.EpiphanyLikeBrowsers,
		browser.LynxLikeBrowsers,

		browserState,
		func(msg string, err error) {
			utils.HandlePanic(cfg.App.Name, msg, err)
		},
	)
}
