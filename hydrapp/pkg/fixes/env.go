//go:build android
// +build android

package fixes

import (
	"os"
	"path/filepath"
)

func PolyfillEnvironment(userHomeDir string) error {
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

	return nil
}
