package utils

import (
	"errors"
	"fmt"
	"os/exec"
	"runtime"

	"github.com/ncruces/zenity"
)

const (
	browserDownloadLink = "https://github.com/pojntfx/hydrapp#which-browsers-are-supported"
)

var (
	ErrNoBrowserOpenMethodFound = errors.New("no method to open a browser found")
)

func HandleNoSupportedBrowserFound(
	appName,

	hydrappBrowserEnv,
	knownBinaries string,
	err error,
) error {
	if err := zenity.Question(
		fmt.Sprintf(`%v requires a supported browser but couldn't find one.

Would you like to download a supported browser or learn more?`, appName),
		zenity.Title(fmt.Sprintf("No supported browser found for %v", appName)),
		zenity.OKLabel("Download"),
		zenity.CancelLabel("Learn more"),
		zenity.Icon(zenity.WarningIcon),
	); err != nil {
		if errors.Is(zenity.ErrCanceled, err) {
			if err := zenity.Question(
				fmt.Sprintf(
					`While searching for a supported browser %v encountered this error:

%v

It tried to find both the preferred the browser binary (set with the HYDRAPP_BROWSER env variable) "%v" and the known binaries:

%v

without success. Would you like to manually configure a browser?`,
					appName,
					err,
					hydrappBrowserEnv,
					knownBinaries,
				),
				zenity.Title(fmt.Sprintf("No supported browser found for %v", appName)),
				zenity.OKLabel("Configure"),
				zenity.CancelLabel("Cancel"),
				zenity.Icon(zenity.InfoIcon),
			); err != nil {
				if errors.Is(zenity.ErrCanceled, err) {
					return err
				} else {
					return fmt.Errorf("could not display dialog: %v", err)
				}
			}

			// TODO: Configure browser manually, set HYDRAPP_BROWSER and HYDRAPP_TYPE and return nil, causing it to re-try

			return nil
		} else {
			return fmt.Errorf("could not display dialog: %v", err)
		}
	}

	switch runtime.GOOS {
	case "windows":
		powerShellBinary, err := exec.LookPath("pwsh.exe")
		if err != nil {
			powerShellBinary = "powershell.exe"
		}

		if output, err := exec.Command(powerShellBinary, `Start-Process`, browserDownloadLink).CombinedOutput(); err != nil {
			return fmt.Errorf("could not open browser with output: %s: %v", output, err)
		}

	case "darwin":
		if output, err := exec.Command("open", browserDownloadLink).CombinedOutput(); err != nil {
			return fmt.Errorf("could not open browser with output: %s: %v", output, err)
		}

	default:
		// Open link with `xdg-open`
		if output, err := exec.Command("xdg-open", browserDownloadLink).CombinedOutput(); err != nil {
			return fmt.Errorf("could not open browser with output: %s: %v", output, err)
		}

		// Open link with `xdg-open`
		if xdgOpen, err := exec.LookPath("xdg-open"); err == nil {
			if output, err := exec.Command(xdgOpen, browserDownloadLink).CombinedOutput(); err != nil {
				return fmt.Errorf("could not open browser with output: %s: %v", output, err)
			}
		} else {
			// Open link with `open` (i.e. FreeBSD and other UNIXes)
			open, err := exec.LookPath("open")
			if err != nil {
				return fmt.Errorf("%v: %w", ErrNoBrowserOpenMethodFound, err)
			}

			if output, err := exec.Command(open, browserDownloadLink).CombinedOutput(); err != nil {
				return fmt.Errorf("could not open browser with output: %s: %v", output, err)
			}
		}
	}

	// TODO: Show "ok"-only modal asking the user to press ok once they have installed the browser, return nil causing it to re-try

	return nil
}
