//go:build selfupdate
// +build selfupdate

package update

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"

	"github.com/blang/semver"
	"github.com/ncruces/zenity"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

func Update(repo string, version string, state *BrowserState) error {
	// Get the latest version
	latest, found, err := selfupdate.DetectLatest(repo)
	if err != nil {
		return err
	}

	// Stop if we are already up to day
	v := semver.MustParse(version)
	if !found || latest.Version.LTE(v) {
		return nil
	}

	// As the user if they want to update
	if cancelled := zenity.Question(
		fmt.Sprintf("A new version (%v) is available, you currently have version %v; do you want to update?", latest, version),
		zenity.Title("Update available"),
		zenity.OKLabel("Update now"),
		zenity.CancelLabel("Remind me later"),
		zenity.Width(320),
	); cancelled != nil {
		return nil
	}

	// Apply the self-update
	self, err := os.Executable()
	if err != nil {
		return err
	}

	if err := selfupdate.UpdateTo(latest.AssetURL, self); err != nil {
		return err
	}

	// Restart self
	cmd := exec.Command(self, os.Args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()

	if state.Cmd != nil && state.Cmd.Process != nil {
		// Windows does not support the `SIGTERM` signal
		if runtime.GOOS == "windows" {
			if output, err := exec.Command("taskkill", "/pid", strconv.Itoa(state.Cmd.Process.Pid)).CombinedOutput(); err != nil {
				return fmt.Errorf("could not close old version: %v: %v", string(output), err)
			}
		} else {
			// We ignore errors here as the old process might already have finished etc.
			_ = state.Cmd.Process.Signal(syscall.SIGTERM)
		}
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Stop old self
	os.Exit(0)

	return nil
}
