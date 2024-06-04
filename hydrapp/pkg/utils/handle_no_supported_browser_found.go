package utils

import (
	"errors"
	"fmt"

	"github.com/ncruces/zenity"
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
				}
			}

			// TODO: Configure browser manually, set HYDRAPP_BROWSER and HYDRAPP_TYPE and return nil, causing it to re-try

			return nil
		} else {
			return err
		}
	}

	// TODO: Open browser download selection screen, show "ok"-only modal asking the user to press ok once they have installed the browser, return nil causing it to re-try

	return nil
}
