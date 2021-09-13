//go:build linux
// +build linux

package fixes

import (
	"os"
)

func init() {
	// Fix Zenity on Wayland when running in Flatpak
	if os.Getenv("XDG_SESSION_TYPE") == "wayland" {
		// Errors are ignored as there is no way to sanely handle errors and there would be no way to recover visually
		_ = os.Setenv("GDK_BACKEND", "wayland")
	}
}
