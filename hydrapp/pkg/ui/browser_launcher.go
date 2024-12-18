package ui

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pojntfx/hydrapp/hydrapp/pkg/utils"
)

const (
	flatpakSpawnCmd  = "flatpak-spawn"
	flatpakSpawnHost = "--host"

	flatpakCmd     = "flatpak"
	flatpakList    = "list"
	flatpakColumns = "--columns=application"
)

var (
	ErrCouldNotFindFlatpakSpawn                              = errors.New("could not find flatpak-spawn command")
	ErrCouldNotListBrowserFlatpaks                           = errors.New("could not list available browser Flatpaks")
	ErrCouldNotCallUnsupportedBrowserHandler                 = errors.New("could not call unsupported browser handler")
	ErrCouldNotLaunchUnknownBrowserType                      = errors.New("could not launch unknown browser type")
	ErrCouldNotGetUserConfigDir                              = errors.New("could not get user config directory")
	ErrCouldNotCreateFirefoxProfile                          = errors.New("could not create Firefox profile")
	ErrCouldNotGetUserHomeDir                                = errors.New("could not get user home directory")
	ErrCouldNotListFilesInFirefoxProfilesDirectory           = errors.New("could not list files in Firefox profiles directory")
	ErrCouldNotFindFirefoxProfileDirectory                   = errors.New("could not find Firefox profile directory")
	ErrCouldNotSetFirefoxProfileDirectoryEnvironmentVariable = errors.New("could not set Firefox profile directory environment variable")
	ErrCouldNotWriteFirefoxPrefsJSFile                       = errors.New("could not write Firefox prefs JS file")
	ErrCouldNotCreateFirefoxChromeDirectory                  = errors.New("could not create Firefox chrome directory")
	ErrCouldNotWriteFirefoxUserChromeCSSFile                 = errors.New("could not write Firefox userChrome CSS file")
	ErrCouldNotCreateEpiphanyProfileDirectory                = errors.New("could not create Epiphany profile directory")
	ErrCouldNotWriteDesktopFile                              = errors.New("could not write desktop file")
	ErrCouldNotSymlinkDesktopFile                            = errors.New("could not symlink desktop file")
	ErrCouldNotFindSupportedBrowser                          = errors.New("could not find a supported browser")
	ErrCouldNotWaitForBrowserLockfileRemoval                 = errors.New("could not wait for browser lockfile removal")
	ErrUnknownBrowserType                                    = errors.New("unknown browser type")
	ErrBrowserLauncherExternalContextCancelled               = errors.New("browser launcher external context cancelled")
	ErrBrowserLauncherFilewatcherContextCancelled            = errors.New("browser launcher file watcher context cancelled")
	ErrBrowserLauncherBrowserCommandBarContextCancelled      = errors.New("browser launcher browser command context cancelled")
)

type Browser struct {
	Name            string
	LinuxBinaries   [][]string
	Flatpaks        [][]string
	WindowsBinaries [][]string
	MacOSBinaries   [][]string
}

func LaunchBrowser(
	ctx context.Context,

	url,
	appName,
	appID,

	browserBinaryOverride,
	browserTypeOverride string,

	chromiumLikeBrowsers Browser,
	firefoxLikeBrowsers Browser,
	epiphanyLikeBrowsers Browser,
	lynxLikeBrowsers Browser,

	state *BrowserState,
	handleNoSupportedBrowserFound func(
		appName,

		hydrappBrowserEnv,
		knownBinaries string,
		err error,
	) error,
) (bool, error) {
	browserBinary := []string{browserBinaryOverride}

	// Process the browser types
	// Order matters; whatever comes first and is discovered first will be used
	rawBrowsers := []Browser{chromiumLikeBrowsers, firefoxLikeBrowsers, epiphanyLikeBrowsers, lynxLikeBrowsers}
	browsers := []Browser{}
	for _, browser := range rawBrowsers {
		// Keep already processed fields
		processedBrowser := Browser{
			browser.Name,
			browser.LinuxBinaries,
			browser.Flatpaks,
			[][]string{},
			[][]string{},
		}

		// Process Windows binaries
		if runtime.GOOS == "windows" {
			for _, suffix := range browser.WindowsBinaries {
				for _, fullPath := range []string{
					filepath.Join(append([]string{os.Getenv("LocalAppData")}, suffix...)...),
					filepath.Join(append([]string{os.Getenv("ProgramFiles")}, suffix...)...),
					filepath.Join(append([]string{os.Getenv("ProgramFiles(x86)")}, suffix...)...),
				} {
					processedBrowser.WindowsBinaries = append(processedBrowser.WindowsBinaries, []string{fullPath})
				}
			}
		}

		// Process macOS binaries
		if runtime.GOOS == "darwin" {
			for _, suffix := range browser.MacOSBinaries {
				processedBrowser.MacOSBinaries = append(processedBrowser.MacOSBinaries, []string{filepath.Join(append([]string{"/Applications"}, suffix...)...)})
			}
		}

		browsers = append(browsers, processedBrowser)
	}

	// Check if we are in Flatpak
	runningInFlatpak := false
	if _, err := exec.LookPath(flatpakSpawnCmd); err == nil {
		runningInFlatpak = true
	}

	// Find browser binary
	browserIsFlatpak := false
	browserIsInSandbox := false
	if browserBinary[0] == "" {
	i:
		for _, browser := range browsers {
			// Find native browser
			for _, binary := range browser.LinuxBinaries {
				// Find supported native browser in native install or in Flatpak sandbox
				if _, err := exec.LookPath(binary[0]); err == nil {
					browserBinary = []string{binary[0]}

					if runningInFlatpak {
						browserIsInSandbox = true
					}

					break i
				}
			}

			// Find Flatpak browser
			if _, err := exec.LookPath(flatpakCmd); err == nil {
				for _, flatpak := range browser.Flatpaks {
					if !runningInFlatpak {
						// Find supported Flatpak browser in native install
						apps, err := exec.CommandContext(ctx, flatpakCmd, flatpakList, flatpakColumns).CombinedOutput()
						if err != nil {
							return false, errors.Join(ErrCouldNotListBrowserFlatpaks, err)
						}

						if strings.Contains(string(apps), flatpak[0]) {
							browserBinary = []string{flatpak[0]}
							browserIsFlatpak = true

							break i
						}
					}
				}
			}

			// Find Windows browser
			if runtime.GOOS == "windows" {
				for _, binary := range browser.WindowsBinaries {
					if _, err := os.Stat(binary[0]); err == nil {
						browserBinary = []string{binary[0]}

						break i
					}
				}
			}

			// Find macOS browser
			if runtime.GOOS == "darwin" {
				for _, binary := range browser.MacOSBinaries {
					if _, err := os.Stat(binary[0]); err == nil {
						browserBinary = []string{binary[0]}

						break i
					}
				}
			}
		}
	}

	if browserBinary[0] == "" && runningInFlatpak {
	k:
		for _, browser := range browsers {
			// Find native browser
			for _, binary := range browser.LinuxBinaries {
				// Find supported native browser on host from Flatpak
				if err := exec.CommandContext(ctx, flatpakSpawnCmd, flatpakSpawnHost, "which", binary[0]).Run(); err == nil {
					browserBinary = []string{binary[0]}

					break k
				}
			}

			// Find Flatpak browser
			if _, err := exec.LookPath(flatpakCmd); err == nil {
				for _, flatpak := range browser.Flatpaks {
					// Find supported Flatpak browser on host from Flatpak
					apps, err := exec.CommandContext(ctx, flatpakSpawnCmd, flatpakSpawnHost, flatpakCmd, flatpakList, flatpakColumns).CombinedOutput()
					if err != nil {
						return false, errors.Join(ErrCouldNotListBrowserFlatpaks, err)
					}

					if strings.Contains(string(apps), flatpak[0]) {
						browserBinary = []string{flatpak[0]}
						browserIsFlatpak = true

						break k
					}
				}
			}
		}
	}

	// Ask for configuration if browser binary could not be found
	if browserBinary[0] == "" {
		if err := handleNoSupportedBrowserFound(
			appName,

			fmt.Sprintf("%v", browserBinary),
			fmt.Sprintf("%v", browsers),

			ErrCouldNotFindSupportedBrowser,
		); err != nil {
			return false, errors.Join(ErrCouldNotCallUnsupportedBrowserHandler, err)
		}

		// Retry if configuration was successful
		return true, nil
	}

	// Find browser type
	if browserTypeOverride == "" {
	j:
		for _, browser := range browsers {
			for _, binary := range append(
				append(
					append(
						browser.LinuxBinaries,
						browser.Flatpaks...,
					),
					browser.WindowsBinaries...,
				),
				browser.MacOSBinaries...,
			) {
				if browserBinary[0] == binary[0] {
					browserTypeOverride = browser.Name

					break j
				}
			}
		}
	}

	// Add `flatpak-run` prefix if browser is Flatpak
	if browserIsFlatpak {
		browserBinary = append([]string{flatpakCmd, "run", "--filesystem=home", "--socket=wayland"}, browserBinary...) // These Flatpak flags are required for Wayland support under Firefox and profile support under Epiphany
	}

	// Add `flatpak-spawn` prefix if running in Flatpak, but the browser is not in the sandbox
	if runningInFlatpak && !browserIsInSandbox {
		browserBinary = append([]string{flatpakSpawnCmd, flatpakSpawnHost}, browserBinary...)
	}

	// Abort if browser type could not be found
	if browserTypeOverride == "" {
		return false, errors.Join(ErrCouldNotLaunchUnknownBrowserType, fmt.Errorf("tried to launch preferred browser type (set with the HYDRAPP_TYPE environment variable) \"%v\" and known types \"%v\"", browserTypeOverride, browsers))
	}

	switch browserTypeOverride {
	// Launch Chromium-like browser
	case BrowserTypeChromium:
		// Create a profile for the app
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return false, errors.Join(ErrCouldNotGetUserConfigDir, err)
		}
		userDataDir := filepath.Join(userConfigDir, appID)

		execArgs := []string{
			"--name=" + appName,
			"--class=" + appName,
			"--user-data-dir=" + userDataDir,
			"--no-first-run",
			"--no-default-browser-check",
			"--app=" + url,
		}

		// If we are on Linux, in a Flatpak sandbox, are running a browser that itself is in the sandbox and
		// are on Wayland, enable the Ozone platform abstraction layer for Chromium to support launching Chromium in Wayland
		if runtime.GOOS == "linux" && runningInFlatpak && browserIsInSandbox && os.Getenv("XDG_SESSION_TYPE") == "wayland" {
			execArgs = append(execArgs, "--ozone-platform-hint=auto")
		}

		// Create the browser instance
		execLine := append(
			browserBinary,
			append(
				execArgs,
				os.Args[1:]...,
			)...,
		)

		state.Cmd = exec.CommandContext(
			ctx,
			execLine[0],
			execLine[1:]...,
		)

		// Use system stdout, stderr and stdin
		state.Cmd.Stdout = os.Stdout
		state.Cmd.Stderr = os.Stderr
		state.Cmd.Stdin = os.Stdin

		// Start the browser
		if err := state.Cmd.Run(); err != nil {
			return false, errors.Join(ErrCouldNotOpenBrowser, err)
		}

		// Wait till lock for browser has been removed
		watch, close, err := utils.SetupFileWatcher(filepath.Join(userDataDir, "SingletonSocket"), false)
		defer close()

		if err != nil {
			return false, errors.Join(ErrCouldNotWaitForBrowserLockfileRemoval, err)
		}

		if err := watch(); err != nil {
			return false, errors.Join(ErrCouldNotWaitForBrowserLockfileRemoval, err)
		}

		// Launch Firefox-like browser
	case BrowserTypeFirefox:
		// Create a profile for the app
		execLine := append(
			browserBinary,
			[]string{
				"--createprofile",
				appID,
			}...,
		)

		if output, err := exec.CommandContext(
			ctx,
			execLine[0],
			execLine[1:]...,
		).CombinedOutput(); err != nil {
			err := fmt.Errorf("could not create Firefox profile with output: %s: %v", output, err)

			return false, errors.Join(ErrCouldNotCreateFirefoxProfile, err)
		}

		// Get the user's home directory in which the profiles can be found
		home, err := os.UserHomeDir()
		if err != nil {
			return false, errors.Join(ErrCouldNotGetUserHomeDir, err)
		}

		// Get the profile's directory
		firefoxDir := filepath.Join(home, ".mozilla", "firefox")
		if _, err := os.Stat(firefoxDir); err != nil {
			// Fall back to non-standard (e.g. snap) profile directories
			firefoxDir = filepath.Join(home, "snap", "firefox", "common", ".mozilla", "firefox")
		}

		if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
			userConfigDir, err := os.UserConfigDir()
			if err != nil {
				return false, errors.Join(ErrCouldNotGetUserConfigDir, err)
			}

			if runtime.GOOS == "windows" {
				firefoxDir = filepath.Join(userConfigDir, "Mozilla", "Firefox", "Profiles")
			} else {
				firefoxDir = filepath.Join(userConfigDir, "Firefox", "Profiles")
			}
		}

		filesInFirefoxDir, err := os.ReadDir(firefoxDir)
		if err != nil {
			return false, errors.Join(ErrCouldNotListFilesInFirefoxProfilesDirectory, err)
		}

		profileSuffix := ""
		for _, file := range filesInFirefoxDir {
			if strings.HasSuffix(file.Name(), appID) {
				profileSuffix = file.Name()

				break
			}
		}

		if profileSuffix == "" {
			return false, ErrCouldNotFindFirefoxProfileDirectory
		}

		profileDir := filepath.Join(firefoxDir, profileSuffix)
		if err := os.Setenv("PROFILE_DIR", profileDir); err != nil {
			return false, errors.Join(ErrCouldNotSetFirefoxProfileDirectoryEnvironmentVariable, err)
		}

		if err := os.WriteFile(filepath.Join(profileDir, "prefs.js"), []byte(prefsJSContent), 0664); err != nil {
			return false, errors.Join(ErrCouldNotWriteFirefoxPrefsJSFile, err)
		}

		chromeDir := filepath.Join(profileDir, "chrome")
		if err := os.MkdirAll(chromeDir, 0755); err != nil {
			return false, errors.Join(ErrCouldNotCreateFirefoxChromeDirectory, err)
		}

		if err := os.WriteFile(filepath.Join(chromeDir, "userChrome.css"), []byte(userChromeCSSContent), 0664); err != nil {
			return false, errors.Join(ErrCouldNotWriteFirefoxUserChromeCSSFile, err)
		}

		// Create the browser instance
		execLine = append(
			browserBinary,
			append(
				[]string{
					"--name=" + appName,
					"--class=" + appName,
					"--new-window",
					"-P",
					appID,
					url,
				},
				os.Args[1:]...,
			)...,
		)

		state.Cmd = exec.CommandContext(
			ctx,
			execLine[0],
			execLine[1:]...,
		)

		// Use system stdout, stderr and stdin
		state.Cmd.Stdout = os.Stdout
		state.Cmd.Stderr = os.Stderr
		state.Cmd.Stdin = os.Stdin

		// Wait till lock for browser has been removed
		// For Firefox, we're waiting for a lockfile modification rather than a removal, so we need to start
		// listening for write events before starting so that creating a second instance doesn't cause it to exit immediately
		watch, close, err := utils.SetupFileWatcher(filepath.Join(profileDir, "storage.sqlite-journal"), true)
		defer close()

		if err != nil {
			return false, errors.Join(ErrCouldNotWaitForBrowserLockfileRemoval, err)
		}

		{
			var errs error

			// We use context.Background here so that we don't confuse a `ctx` cancellation, after which we should return, with a successful download
			// We select beteween `filewatcherCtx` and `ctx` on all code paths so that we don't leak it
			filewatcherCtx, cancelFilewatcherContext := context.WithCancel(context.Background())
			defer cancelFilewatcherContext()

			go func() {
				defer cancelFilewatcherContext()

				if err := watch(); err != nil {
					errs = errors.Join(errs, errors.Join(ErrCouldNotWaitForBrowserLockfileRemoval, err))

					return
				}
			}()

			// We use context.Background here so that we don't confuse a `ctx` cancellation, after which we should return, with a successful progress bar
			// We select beteween `browserCtx` and `ctx` on all code paths so that we don't leak it
			browserCtx, cancelBrowserCtx := context.WithCancel(context.Background())
			defer cancelBrowserCtx()

			go func() {
				defer cancelBrowserCtx()

				// Start the browser
				if err := state.Cmd.Run(); err != nil {
					errs = errors.Join(errs, errors.Join(ErrCouldNotOpenBrowser, err))

					return
				}
			}()

			select {
			case <-ctx.Done():
				if err := ctx.Err(); err != context.Canceled || errs != nil {
					return false, errors.Join(ErrBrowserLauncherExternalContextCancelled, errs, err)
				}

			case <-filewatcherCtx.Done():
				if err := filewatcherCtx.Err(); err != context.Canceled || errs != nil {
					return false, errors.Join(ErrBrowserLauncherFilewatcherContextCancelled, errs, err)
				}

				select {
				case <-ctx.Done():
					if err := ctx.Err(); err != context.Canceled || errs != nil {
						return false, errors.Join(ErrBrowserLauncherExternalContextCancelled, errs, err)
					}

				case <-browserCtx.Done():
					if err := browserCtx.Err(); err != context.Canceled || errs != nil {
						return false, errors.Join(ErrBrowserLauncherBrowserCommandBarContextCancelled, errs, err)
					}
				}

			case <-browserCtx.Done():
				if err := browserCtx.Err(); err != context.Canceled || errs != nil {
					return false, errors.Join(ErrBrowserLauncherBrowserCommandBarContextCancelled, errs, err)
				}

				select {
				case <-ctx.Done():
					if err := ctx.Err(); err != context.Canceled || errs != nil {
						return false, errors.Join(ErrBrowserLauncherExternalContextCancelled, errs, err)
					}

				case <-filewatcherCtx.Done():
					if err := filewatcherCtx.Err(); err != context.Canceled || errs != nil {
						return false, errors.Join(ErrBrowserLauncherFilewatcherContextCancelled, errs, err)
					}
				}
			}
		}

		// Launch Epiphany-like browser
	case BrowserTypeEpiphany:
		// Get the user's data directory in which the profiles should be created
		dataHomeDir := os.Getenv("XDG_DATA_HOME")
		if strings.TrimSpace(dataHomeDir) == "" {
			userHomeDir, err := os.UserHomeDir()
			if err != nil {
				return false, errors.Join(ErrCouldNotGetUserHomeDir, err)
			}

			dataHomeDir = filepath.Join(userHomeDir, ".local", "share")
		}

		// Create the profile directory
		epiphanyID := "org.gnome.Epiphany.WebApp_" + appID

		applicationsDir := filepath.Join(dataHomeDir, "applications")
		if err := os.MkdirAll(applicationsDir, 0755); err != nil {
			return false, errors.Join(ErrCouldNotCreateEpiphanyProfileDirectory, err)
		}

		profileDir := filepath.Join(dataHomeDir, epiphanyID)
		if err := os.MkdirAll(filepath.Join(profileDir, ".app"), 0755); err != nil {
			panic(err)
		}
		defer os.RemoveAll(profileDir)

		// Create the .desktop file
		desktopFilePath := filepath.Join(applicationsDir, epiphanyID+".desktop")
		if err := os.WriteFile(
			desktopFilePath,
			[]byte(fmt.Sprintf(
				epiphanyDesktopFileTemplate,
				browserBinary[0],
				profileDir,
				url,
				epiphanyID,
				appName,
				appID,
			)),
			0664,
		); err != nil {
			return false, errors.Join(ErrCouldNotWriteDesktopFile, err)
		}
		defer os.RemoveAll(desktopFilePath)

		// Symlink .desktop file to expected XDG directories
		xdgApplicationsDir := filepath.Join(dataHomeDir, "xdg-desktop-portal", "applications")
		if err := os.MkdirAll(xdgApplicationsDir, 0755); err != nil {
			return false, errors.Join(ErrCouldNotSymlinkDesktopFile, err)
		}

		xdgApplicationsDesktopFilePath := filepath.Join(xdgApplicationsDir, epiphanyID+".desktop")
		if err := os.Symlink(desktopFilePath, xdgApplicationsDesktopFilePath); err != nil {
			return false, errors.Join(ErrCouldNotSymlinkDesktopFile, err)
		}
		defer os.RemoveAll(xdgApplicationsDesktopFilePath)

		// Create the browser instance
		execLine := append(
			browserBinary,
			append(
				[]string{
					"--new-window",
					"--application-mode",
					"--profile=" + profileDir,
					url,
				},
				os.Args[1:]...,
			)...,
		)

		state.Cmd = exec.CommandContext(
			ctx,
			execLine[0],
			execLine[1:]...,
		)

		// Use system stdout, stderr and stdin
		state.Cmd.Stdout = os.Stdout
		state.Cmd.Stderr = os.Stderr
		state.Cmd.Stdin = os.Stdin

		// Start the browser
		if err := state.Cmd.Run(); err != nil {
			return false, errors.Join(ErrCouldNotOpenBrowser, err)
		}

		// Wait till lock for browser has been removed
		watch, close, err := utils.SetupFileWatcher(filepath.Join(profileDir, "ephy-history.db-wal"), false)
		defer close()

		if err != nil {
			return false, errors.Join(ErrCouldNotWaitForBrowserLockfileRemoval, err)
		}

		if err := watch(); err != nil {
			return false, errors.Join(ErrCouldNotWaitForBrowserLockfileRemoval, err)
		}

		// Launch Lynx-like browser
	case BrowserTypeLynx:
		// Create the browser instance
		execLine := append(
			browserBinary,
			append(
				[]string{
					"--nopause",
					"--accept_all_cookies",
					url,
				},
				os.Args[1:]...,
			)...,
		)

		state.Cmd = exec.CommandContext(
			ctx,
			execLine[0],
			execLine[1:]...,
		)

		// Use system stdout, stderr and stdin
		state.Cmd.Stdout = os.Stdout
		state.Cmd.Stderr = os.Stderr
		state.Cmd.Stdin = os.Stdin

		// Start the browser
		if err := state.Cmd.Run(); err != nil {
			return false, errors.Join(ErrCouldNotOpenBrowser, err)
		}

		// No need to wait till lock for browser has been removed in the case of Lynx since there are no profiles

	// Launch dummy browser
	case BrowserTypeDummy:
		select {}
	default:
		if err := handleNoSupportedBrowserFound(
			appName,

			fmt.Sprintf("%v", browserBinary),
			fmt.Sprintf("%v", browsers),

			ErrUnknownBrowserType,
		); err != nil {
			return false, errors.Join(ErrCouldNotCallUnsupportedBrowserHandler, err)
		}

		// Retry if configuration was successful
		return true, nil
	}

	return false, nil
}
