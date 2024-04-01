package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

const (
	EnvBackendLaddr  = "HYDRAPP_BACKEND_LADDR"
	EnvFrontendLaddr = "HYDRAPP_FRONTEND_LADDR"
	EnvBrowser       = "HYDRAPP_BROWSER"
	EnvType          = "HYDRAPP_TYPE"
	EnvSelfupdate    = "HYDRAPP_SELFUPDATE"
)

func PolyfillEnvironment(userHomeDir string) error {
	switch runtime.GOOS {
	case "android":
		if _, exists := os.LookupEnv("XDG_CACHE_HOME"); !exists {
			if err := os.Setenv("XDG_CACHE_HOME", filepath.Join(userHomeDir, ".cache")); err != nil {
				return err
			}
		}

		if _, exists := os.LookupEnv("XDG_CONFIG_HOME"); !exists {
			if err := os.Setenv("XDG_CONFIG_HOME", filepath.Join(userHomeDir, ".config")); err != nil {
				return err
			}
		}

		if _, exists := os.LookupEnv("HOME"); !exists {
			if err := os.Setenv("HOME", userHomeDir); err != nil {
				return err
			}
		}
	}

	return nil
}
