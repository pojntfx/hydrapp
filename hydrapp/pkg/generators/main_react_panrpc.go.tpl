//go:build !android
// +build !android

package main

import (
	"bytes"
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"{{ .GoMod }}/pkg/backend"
	"{{ .GoMod }}/pkg/frontend"
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
		ui.HandlePanic("App", errors.Join(fmt.Errorf("could not parse config file"), err))

		return
	}

	// Apply the self-update
	browserState := &ui.BrowserState{}
	go func() {
		if err := ui.SelfUpdate(
			ctx,

			cfg,
			browserState,
		); err != nil {
			ui.HandlePanic(cfg.App.Name, err)
		}
	}()

	// Start the backend
	backendURL, stopBackend, err := backend.StartServer(ctx, os.Getenv(ui.EnvBackendLaddr), time.Second*10, true)
	if err != nil {
		ui.HandlePanic(cfg.App.Name, errors.Join(fmt.Errorf("could not start backend"), err))
	}
	defer stopBackend()

	log.Println("Backend URL:", backendURL)

	// Start the frontend
	frontendURL, stopFrontend, err := frontend.StartServer(ctx, os.Getenv(ui.EnvFrontendLaddr), backendURL, true)
	if err != nil {
		ui.HandlePanic(cfg.App.Name, errors.Join(fmt.Errorf("could not start frontend"), err))
	}
	defer stopFrontend()

	log.Println("Frontend URL:", frontendURL)

	for {
		retry, err := ui.LaunchBrowser(
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
			ui.ConfigureBrowser,
		)

		if err != nil {
			ui.HandlePanic(cfg.App.Name, err)
		}

		if !retry {
			return
		}
	}
}
