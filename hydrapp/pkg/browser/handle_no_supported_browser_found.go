package browser

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/ncruces/zenity"
	"github.com/pojntfx/hydrapp/hydrapp/pkg/utils"
)

const (
	browserDownloadLink = "https://github.com/pojntfx/hydrapp#which-browsers-are-supported"

	browserDescriptionChromium = "Chromium-like (Chrome, Edge, Brave etc.)"
	browserDescriptionFirefox  = "Firefox"
	browserDescriptionEpiphany = "GNOME Web/Epiphany"
	browserDescriptionLynx     = "Lynx"

	browserDescriptionDummy = "Dummy/no browser"
)

var (
	ErrNoBrowserOpenMethodFound      = errors.New("no method to open a browser found")
	ErrBrowserConfigurationCancelled = errors.New("browser configuration cancelled")
)

func HandleNoSupportedBrowserFound(
	appName,

	hydrappBrowserEnv,
	knownBinaries string,
	err error,
) error {
	if e := zenity.Question(
		fmt.Sprintf(`%v requires a supported browser but couldn't find one.

Would you like to download a supported browser or learn more?`, appName),
		zenity.Title(fmt.Sprintf("No supported browser found for %v", appName)),
		zenity.OKLabel("Download"),
		zenity.CancelLabel("Learn more"),
		zenity.Icon(zenity.WarningIcon),
	); e != nil {
		if errors.Is(zenity.ErrCanceled, e) {
			if err := zenity.Question(
				fmt.Sprintf(
					`While searching for a supported browser %v encountered this error:

%v

It tried to find both the preferred the browser binary (set with the HYDRAPP_BROWSER environment variable) "%v" and the known binaries:

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
					return ErrBrowserConfigurationCancelled
				}

				return fmt.Errorf("could not display dialog: %v", err)
			}

			browserDescription, err := zenity.List(
				"Select your browser type",
				[]string{
					browserDescriptionChromium,
					browserDescriptionFirefox,
					browserDescriptionEpiphany,
					browserDescriptionLynx,

					browserDescriptionDummy,
				},
				zenity.Title(fmt.Sprintf("Browser type configuration for %v", appName)),
				zenity.OKLabel("Continue"),
			)
			if err != nil {
				if errors.Is(zenity.ErrCanceled, err) {
					return ErrBrowserConfigurationCancelled
				}

				return fmt.Errorf("could not display dialog: %v", err)
			}

			switch browserDescription {
			case browserDescriptionChromium:
				if err := os.Setenv(utils.EnvType, browserTypeChromium); err != nil {
					return fmt.Errorf("could not set environment variable: %v", err)
				}

			case browserDescriptionFirefox:
				if err := os.Setenv(utils.EnvType, browserTypeFirefox); err != nil {
					return fmt.Errorf("could not set environment variable: %v", err)
				}

			case browserDescriptionEpiphany:
				if err := os.Setenv(utils.EnvType, browserTypeEpiphany); err != nil {
					return fmt.Errorf("could not set environment variable: %v", err)
				}

			case browserDescriptionLynx:
				if err := os.Setenv(utils.EnvType, browserTypeLynx); err != nil {
					return fmt.Errorf("could not set environment variable: %v", err)
				}

			// No need to check extra options here since it's a radio select and only valid options can be returned
			default:
				if err := os.Setenv(utils.EnvType, browserTypeDummy); err != nil {
					return fmt.Errorf("could not set environment variable: %v", err)
				}
			}

			browserLocation, err := zenity.Entry(
				"Browser binary location or command:",
				zenity.Title(fmt.Sprintf("Browser location configuration for %v", appName)),
				zenity.OKLabel("Continue"),
			)
			if err != nil {
				if errors.Is(zenity.ErrCanceled, err) {
					return ErrBrowserConfigurationCancelled
				}

				return fmt.Errorf("could not display dialog: %v", err)
			}

			if err := os.Setenv(utils.EnvBrowser, browserLocation); err != nil {
				return fmt.Errorf("could not set environment variable: %v", err)
			}

			return nil
		} else {
			return fmt.Errorf("could not display dialog: %v", e)
		}
	}

	switch runtime.GOOS {
	case "windows":
		powerShellBinary, err := exec.LookPath("pwsh.exe")
		if err != nil {
			powerShellBinary = "powershell.exe"
		}

		if output, err := exec.Command(powerShellBinary, `-Command`, fmt.Sprintf(`Start-Process %v`, browserDownloadLink)).CombinedOutput(); err != nil {
			return fmt.Errorf("could not open browser with output: %s: %v", output, err)
		}

	case "darwin":
		if output, err := exec.Command("open", browserDownloadLink).CombinedOutput(); err != nil {
			return fmt.Errorf("could not open browser with output: %s: %v", output, err)
		}

	default:
		// Open link with `xdg-open` (we need to detach because `xdg-open` may block, unlike the Windows and macOS equivalents)
		if xdgOpen, err := exec.LookPath("xdg-open"); err == nil {
			cmd := exec.Command(xdgOpen, browserDownloadLink)

			var output bytes.Buffer
			cmd.Stdout = &output
			cmd.Stderr = &output

			if err := cmd.Run(); err != nil {
				return fmt.Errorf("could not open browser with output: %s: %v", output.String(), err)
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

	if err := zenity.Info(
		"Continue once you have downloaded and installed a supported browser",
		zenity.Title(fmt.Sprintf("%v is waiting for browser installation", appName)),
		zenity.OKLabel("Continue"),
	); err != nil {
		if errors.Is(zenity.ErrCanceled, err) {
			return ErrBrowserConfigurationCancelled
		}

		return fmt.Errorf("could not display dialog: %v", err)
	}

	return nil
}
