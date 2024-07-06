package ui

import (
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
	ErrCouldNotFindSupportedBrowser                          = errors.New("could not find a supported browser")
	ErrCouldNotWaitForBrowserLockfileRemoval                 = errors.New("could not wait for browser lockfile removal")
)

type Browser struct {
	Name            string
	LinuxBinaries   [][]string
	Flatpaks        [][]string
	WindowsBinaries [][]string
	MacOSBinaries   [][]string
}

func LaunchBrowser(
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
	handlePanic func(msg string, err error),
	handleNoSupportedBrowserFound func(
		appName,

		hydrappBrowserEnv,
		knownBinaries string,
		err error,
	) error,
) bool {
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
	if browserBinary[0] == "" {
	i:
		for _, browser := range browsers {
			// Find native browser
			for _, binary := range browser.LinuxBinaries {
				if runningInFlatpak {
					// Find supported browser from Flatpak
					if err := exec.Command(flatpakSpawnCmd, flatpakSpawnHost, "which", binary[0]).Run(); err == nil {
						browserBinary = []string{binary[0]}

						break i
					}
				} else {
					// Find supported browser in native install
					if _, err := exec.LookPath(binary[0]); err == nil {
						browserBinary = []string{binary[0]}

						break i
					}
				}
			}

			// Find Flatpak browser
			if _, err := exec.LookPath(flatpakCmd); err == nil {
				for _, flatpak := range browser.Flatpaks {
					if runningInFlatpak {
						// Find supported browser from Flatpak
						apps, err := exec.Command(flatpakSpawnCmd, flatpakSpawnHost, flatpakCmd, flatpakList, flatpakColumns).CombinedOutput()
						if err != nil {
							handlePanic(ErrCouldNotListBrowserFlatpaks.Error(), errors.Join(ErrCouldNotListBrowserFlatpaks, err))
						}

						if strings.Contains(string(apps), flatpak[0]) {
							browserBinary = []string{flatpak[0]}
							browserIsFlatpak = true

							break i
						}
					} else {
						// Find supported browser in native install
						apps, err := exec.Command(flatpakCmd, flatpakList, flatpakColumns).CombinedOutput()
						if err != nil {
							handlePanic(ErrCouldNotListBrowserFlatpaks.Error(), errors.Join(ErrCouldNotListBrowserFlatpaks, err))
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

	// Ask for configuration if browser binary could not be found
	if browserBinary[0] == "" {
		if err := handleNoSupportedBrowserFound(
			appName,

			fmt.Sprintf("%v", browserBinary),
			fmt.Sprintf("%v", browsers),

			ErrCouldNotFindSupportedBrowser,
		); err != nil {
			handlePanic(ErrCouldNotCallUnsupportedBrowserHandler.Error(), errors.Join(ErrCouldNotCallUnsupportedBrowserHandler, err))
		}

		// Retry if configuration was successful
		return true
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

	// Add `flatpak-spawn` prefix if running in Flatpak
	if runningInFlatpak {
		browserBinary = append([]string{flatpakSpawnCmd, flatpakSpawnHost}, browserBinary...)
	}

	// Abort if browser type could not be found
	if browserTypeOverride == "" {
		handlePanic(ErrCouldNotLaunchUnknownBrowserType.Error(), errors.Join(ErrCouldNotLaunchUnknownBrowserType, fmt.Errorf("tried to launch preferred browser type (set with the HYDRAPP_TYPE environment variable) \"%v\" and known types \"%v\"", browserTypeOverride, browsers)))
	}

	switch browserTypeOverride {
	// Launch Chromium-like browser
	case BrowserTypeChromium:
		// Create a profile for the app
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			handlePanic(ErrCouldNotGetUserConfigDir.Error(), errors.Join(ErrCouldNotGetUserConfigDir, err))
		}
		userDataDir := filepath.Join(userConfigDir, appID)

		// Create the browser instance
		execLine := append(
			browserBinary,
			append(
				[]string{
					"--name=" + appName,
					"--class=" + appName,
					"--user-data-dir=" + userDataDir,
					"--no-first-run",
					"--no-default-browser-check",
					"--app=" + url,
				},
				os.Args[1:]...,
			)...,
		)

		state.Cmd = exec.Command(
			execLine[0],
			execLine[1:]...,
		)

		// Use system stdout, stderr and stdin
		state.Cmd.Stdout = os.Stdout
		state.Cmd.Stderr = os.Stderr
		state.Cmd.Stdin = os.Stdin

		// Start the browser
		if err := state.Cmd.Run(); err != nil {
			handlePanic(ErrCouldNotOpenBrowser.Error(), errors.Join(ErrCouldNotOpenBrowser, err))
		}

		// Wait till lock for browser has been removed
		if err := utils.WaitForFileRemoval(filepath.Join(userDataDir, "SingletonSocket")); err != nil {
			handlePanic(ErrCouldNotWaitForBrowserLockfileRemoval.Error(), errors.Join(ErrCouldNotWaitForBrowserLockfileRemoval, err))
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

		if output, err := exec.Command(
			execLine[0],
			execLine[1:]...,
		).CombinedOutput(); err != nil {
			err := fmt.Errorf("could not create Firefox profile with output: %s: %v", output, err)

			handlePanic(ErrCouldNotCreateFirefoxProfile.Error(), errors.Join(ErrCouldNotCreateFirefoxProfile, err))
		}

		// Get the user's home directory in which the profiles can be found
		home, err := os.UserHomeDir()
		if err != nil {
			handlePanic(ErrCouldNotGetUserHomeDir.Error(), errors.Join(ErrCouldNotGetUserHomeDir, err))
		}

		// Get the profile's directory
		firefoxDir := filepath.Join(home, ".mozilla", "firefox")
		if runtime.GOOS == "windows" || runtime.GOOS == "darwin" {
			userConfigDir, err := os.UserConfigDir()
			if err != nil {
				handlePanic(ErrCouldNotGetUserConfigDir.Error(), errors.Join(ErrCouldNotGetUserConfigDir, err))
			}

			if runtime.GOOS == "windows" {
				firefoxDir = filepath.Join(userConfigDir, "Mozilla", "Firefox", "Profiles")
			} else {
				firefoxDir = filepath.Join(userConfigDir, "Firefox", "Profiles")
			}
		}

		filesInFirefoxDir, err := os.ReadDir(firefoxDir)
		if err != nil {
			handlePanic(ErrCouldNotListFilesInFirefoxProfilesDirectory.Error(), errors.Join(ErrCouldNotListFilesInFirefoxProfilesDirectory, err))
		}

		profileSuffix := ""
		for _, file := range filesInFirefoxDir {
			if strings.HasSuffix(file.Name(), appID) {
				profileSuffix = file.Name()

				break
			}
		}

		if profileSuffix == "" {
			handlePanic(ErrCouldNotFindFirefoxProfileDirectory.Error(), ErrCouldNotFindFirefoxProfileDirectory)
		}

		profileDir := filepath.Join(firefoxDir, profileSuffix)
		if err := os.Setenv("PROFILE_DIR", profileDir); err != nil {
			handlePanic(ErrCouldNotSetFirefoxProfileDirectoryEnvironmentVariable.Error(), errors.Join(ErrCouldNotSetFirefoxProfileDirectoryEnvironmentVariable, err))
		}

		if err := os.WriteFile(filepath.Join(profileDir, "prefs.js"), []byte(prefsJSContent), 0664); err != nil {
			handlePanic(ErrCouldNotWriteFirefoxPrefsJSFile.Error(), errors.Join(ErrCouldNotWriteFirefoxPrefsJSFile, err))
		}

		chromeDir := filepath.Join(profileDir, "chrome")
		if err := os.MkdirAll(chromeDir, 0755); err != nil {
			handlePanic(ErrCouldNotCreateFirefoxChromeDirectory.Error(), errors.Join(ErrCouldNotCreateFirefoxChromeDirectory, err))
		}

		if err := os.WriteFile(filepath.Join(chromeDir, "userChrome.css"), []byte(userChromeCSSContent), 0664); err != nil {
			handlePanic(ErrCouldNotWriteFirefoxUserChromeCSSFile.Error(), errors.Join(ErrCouldNotWriteFirefoxUserChromeCSSFile, err))
		}

		// Create the browser instance
		execLine = append(
			browserBinary,
			append(
				[]string{
					"--name=" + appName,
					"--class=" + appName,
					"--new-window",
					"--no-first-run",
					"-P",
					appID,
					url,
				},
				os.Args[1:]...,
			)...,
		)

		state.Cmd = exec.Command(
			execLine[0],
			execLine[1:]...,
		)

		// Use system stdout, stderr and stdin
		state.Cmd.Stdout = os.Stdout
		state.Cmd.Stderr = os.Stderr
		state.Cmd.Stdin = os.Stdin

		// Start the browser
		if err := state.Cmd.Run(); err != nil {
			handlePanic(ErrCouldNotOpenBrowser.Error(), errors.Join(ErrCouldNotOpenBrowser, err))
		}

		// Wait till lock for browser has been removed
		if err := utils.WaitForFileRemoval(filepath.Join(profileDir, "cookies.sqlite-wal")); err != nil {
			handlePanic(ErrCouldNotWaitForBrowserLockfileRemoval.Error(), errors.Join(ErrCouldNotWaitForBrowserLockfileRemoval, err))
		}

		// Launch Epiphany-like browser
	case BrowserTypeEpiphany:
		// Get the user's home directory in which the profiles should be created
		home, err := os.UserHomeDir()
		if err != nil {
			handlePanic(ErrCouldNotGetUserHomeDir.Error(), errors.Join(ErrCouldNotGetUserHomeDir, err))
		}

		// Create the profile directory
		epiphanyID := "org.gnome.Epiphany.WebApp_" + appID
		profileDir := filepath.Join(home, ".local", "share", epiphanyID)

		if err := os.MkdirAll(filepath.Join(profileDir, ".app"), 0755); err != nil {
			handlePanic(ErrCouldNotCreateEpiphanyProfileDirectory.Error(), errors.Join(ErrCouldNotCreateEpiphanyProfileDirectory, err))
		}

		// Create the .desktop file
		if err := os.WriteFile(
			filepath.Join(profileDir, epiphanyID+".desktop"),
			[]byte(fmt.Sprintf(
				epiphanyDesktopFileTemplate,
				appName,
				appName,
				profileDir,
				url,
				appName,
				epiphanyID,
			)),
			0664); err != nil {
			handlePanic(ErrCouldNotWriteDesktopFile.Error(), errors.Join(ErrCouldNotWriteDesktopFile, err))
		}

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

		state.Cmd = exec.Command(
			execLine[0],
			execLine[1:]...,
		)

		// Use system stdout, stderr and stdin
		state.Cmd.Stdout = os.Stdout
		state.Cmd.Stderr = os.Stderr
		state.Cmd.Stdin = os.Stdin

		// Start the browser
		if err := state.Cmd.Run(); err != nil {
			handlePanic(ErrCouldNotOpenBrowser.Error(), errors.Join(ErrCouldNotOpenBrowser, err))
		}

		// Wait till lock for browser has been removed
		if err := utils.WaitForFileRemoval(filepath.Join(profileDir, "ephy-history.db-wal")); err != nil {
			handlePanic(ErrCouldNotWaitForBrowserLockfileRemoval.Error(), errors.Join(ErrCouldNotWaitForBrowserLockfileRemoval, err))
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

		state.Cmd = exec.Command(
			execLine[0],
			execLine[1:]...,
		)

		// Use system stdout, stderr and stdin
		state.Cmd.Stdout = os.Stdout
		state.Cmd.Stderr = os.Stderr
		state.Cmd.Stdin = os.Stdin

		// Start the browser
		if err := state.Cmd.Run(); err != nil {
			handlePanic(ErrCouldNotOpenBrowser.Error(), errors.Join(ErrCouldNotOpenBrowser, err))
		}
	case BrowserTypeDummy:
		select {}
	}

	return false
}
