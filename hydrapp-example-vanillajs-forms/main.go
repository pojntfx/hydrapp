//go:build !android
// +build !android

package main

import (
	"bytes"
	"context"
	_ "embed"
	"log"
	"os"

	"github.com/pojntfx/hydrapp/hydrapp-example-vanillajs-forms/pkg/frontend"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/config"
	_ "github.com/pojntfx/hydrapp/hydrapp/pkg/fixes"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/ui"
)

//go:embed hydrapp.yaml
var configFile []byte

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Parse(bytes.NewBuffer(configFile))
	if err != nil {
		ui.HandlePanic("App", "could not parse config file", err)

		return
	}

	// Apply the self-update
	browserState := &ui.BrowserState{}
	go ui.SelfUpdate(
		ctx,

		cfg,
		browserState,
		ui.HandlePanic,
	)

	// Start the frontend
	frontendURL, stopFrontend, err := frontend.StartServer(ctx, os.Getenv(ui.EnvFrontendLaddr), true)
	if err != nil {
		ui.HandlePanic(cfg.App.Name, "could not start frontend", err)
	}
	defer stopFrontend()

	log.Println("Frontend URL:", frontendURL)

	for {
		if !ui.LaunchBrowser(
			frontendURL,
			cfg.App.Name,
			cfg.App.ID,

			os.Getenv(ui.EnvBrowser),
			os.Getenv(ui.EnvType),

			ui.ChromiumLikeBrowsers,
			ui.FirefoxLikeBrowsers,
			ui.EpiphanyLikeBrowsers,
			ui.LynxLikeBrowsers,

			browserState,
			func(msg string, err error) {
				ui.HandlePanic(cfg.App.Name, msg, err)
			},
			ui.ConfigureBrowser,
		) {
			return
		}
	}
}
