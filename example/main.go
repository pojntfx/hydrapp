//go:build !android
// +build !android

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"github.com/fsnotify/fsnotify"
	"github.com/ncruces/zenity"
	"github.com/pojntfx/hydrapp/example/pkg/backend"
	_ "github.com/pojntfx/hydrapp/example/pkg/fixes"
)

const (
	name = "Hydrapp Example"
	id   = "com.pojtinger.felicitas.hydrapp.example"

	flatpakSpawnCmd  = "flatpak-spawn"
	flatpakSpawnHost = "--host"

	flatpakCmd     = "flatpak"
	flatpakList    = "list"
	flatpakColumns = "--columns=application"

	browserTypeChromium = "chromium"
	browserTypeFirefox  = "firefox"
	browserTypeEpiphany = "epiphany"

	osTypeWindows = "windows"
	osTypeMacOS   = "darwin"

	firefoxPrefs = `user_pref("toolkit.legacyUserProfileCustomizations.stylesheets", true);
user_pref("browser.tabs.drawInTitlebar", false);`
	firefoxUserChrome = `TabsToolbar {
  visibility: collapse;
}
:root:not([customizing]) #navigator-toolbox:not(:hover):not(:focus-within) {
  max-height: 1px;
  min-height: calc(0px);
  overflow: hidden;
}
#navigator-toolbox::after {
  display: none !important;
}
#main-window[sizemode="maximized"] #content-deck {
  padding-top: 8px;
}`

	epiphanyDesktopFileTemplate = `[Desktop Entry]
Name=%v
Exec=epiphany --name="%v" --class="%v" --new-window --application-mode --profile="%v" "%v"
StartupNotify=true
Terminal=false
Type=Application
Categories=GNOME;GTK;
StartupWMClass=%v
X-Purism-FormFactor=Workstation;Mobile;`
)

type Browser struct {
	Name            string
	LinuxBinaries   [][]string
	Flatpaks        [][]string
	WindowsBinaries [][]string
	MacOSBinaries   [][]string
}

var chromiumLikeBrowsers = Browser{
	Name: browserTypeChromium,
	LinuxBinaries: [][]string{
		{"google-chrome"},
		{"google-chrome-stable"},
		{"google-chrome-beta"},
		{"google-chrome-unstable"},
		{"brave-browser"},
		{"brave-browser-stable"},
		{"brave-browser-beta"},
		{"brave-browser-nightly"},
		{"microsoft-edge"},
		{"microsoft-edge-beta"},
		{"microsoft-edge-dev"},
		{"microsoft-edge-canary"},
		{"ungoogled-chromium"},
		{"chromium-browser"},
		{"chromium"},
	},
	Flatpaks: [][]string{
		{"org.chromium.Chromium"},
		{"com.github.Eloston.UngoogledChromium"},
	},
	WindowsBinaries: [][]string{
		{"Google", "Chrome", "Application", "chrome.exe"},
		{"Google", "Chrome Beta", "Application", "chrome.exe"},
		{"Google", "Chrome SxS", "Application", "chrome.exe"},
		{"BraveSoftware", "Brave-Browser", "Application", "brave.exe"},
		{"BraveSoftware", "Brave-Browser-Beta", "Application", "brave.exe"},
		{"BraveSoftware", "Brave-Browser-Nightly", "Application", "brave.exe"},
		{"Microsoft", "Edge", "Application", "msedge.exe"},
		{"Microsoft", "Edge Beta", "Application", "msedge.exe"},
		{"Microsoft", "Edge Dev", "Application", "msedge.exe"},
		{"Microsoft", "Edge Canary", "Application", "msedge.exe"},
		{"Chromium", "Application", "chrome.exe"},
	},
	MacOSBinaries: [][]string{
		{"Google Chrome.app", "Contents", "MacOS", "Google Chrome"},
		{"Google Chrome Beta.app", "Contents", "MacOS", "Google Chrome Beta"},
		{"Google Chrome Canary.app", "Contents", "MacOS", "Google Chrome Canary"},
		{"Brave Browser.app", "Contents", "MacOS", "Brave Browser"},
		{"Brave Browser Beta.app", "Contents", "MacOS", "Brave Browser Beta"},
		{"Brave Browser Nightly.app", "Contents", "MacOS", "Brave Browser Nightly"},
		{"Microsoft Edge.app", "Contents", "MacOS", "Microsoft Edge"},
		{"Microsoft Edge Beta.app", "Contents", "MacOS", "Microsoft Edge Beta"},
		{"Microsoft Edge Dev.app", "Contents", "MacOS", "Microsoft Edge Dev"},
		{"Microsoft Edge Canary.app", "Contents", "MacOS", "Microsoft Edge Canary"},
		{"Chromium.app", "Contents", "MacOS", "Chromium"},
	},
}

var firefoxLikeBrowsers = Browser{
	Name: browserTypeFirefox,
	LinuxBinaries: [][]string{
		{"firefox"},
		{"firefox-esr"},
	},
	Flatpaks: [][]string{
		{"org.mozilla.firefox"},
	},
	WindowsBinaries: [][]string{
		{"Mozilla Firefox", "firefox.exe"},
		{"Firefox Nightly", "firefox.exe"},
	},
	MacOSBinaries: [][]string{
		{"Firefox.app", "Contents", "MacOS", "firefox"},
		{"Firefox Nightly.app", "Contents", "MacOS", "firefox"},
	},
}

var epiphanyLikeBrowsers = Browser{
	Name: browserTypeEpiphany,
	LinuxBinaries: [][]string{
		{"epiphany"},
	},
	Flatpaks: [][]string{
		{"org.gnome.Epiphany"},
	},
	WindowsBinaries: [][]string{},
	MacOSBinaries:   [][]string{},
}

func main() {
	// Start the integrated webserver server
	url, stop, err := backend.StartServer()
	if err != nil {
		crash("could not start integrated webserver", err)
	}
	defer stop()

	// Use the user-prefered browser binary and type if specified
	browserBinary := []string{os.Getenv("HYDRAPP_BROWSER")}
	browserType := os.Getenv("HYDRAPP_TYPE")

	// Process the browser types
	// Order matters; whatever comes first and is discovered first will be used
	rawBrowsers := []Browser{chromiumLikeBrowsers, firefoxLikeBrowsers, epiphanyLikeBrowsers}
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
		if runtime.GOOS == osTypeWindows {
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
		if runtime.GOOS == osTypeMacOS {
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
							crash("could not list available browser flatpaks", err)
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
							crash("could not list available browser flatpaks", err)
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
			if runtime.GOOS == osTypeWindows {
				for _, binary := range browser.WindowsBinaries {
					if _, err := os.Stat(binary[0]); err == nil {
						browserBinary = []string{binary[0]}

						break i
					}
				}
			}

			// Find macOS browser
			if runtime.GOOS == osTypeMacOS {
				for _, binary := range browser.MacOSBinaries {
					if _, err := os.Stat(binary[0]); err == nil {
						browserBinary = []string{binary[0]}

						break i
					}
				}
			}
		}
	}

	// Abort if browser binary could not be found
	if browserBinary[0] == "" {
		crash("could not find a supported browser", fmt.Errorf("tried to launch preferred browser binary (set with the HYDRAPP_BROWSER env variable) \"%v\" and known binaries \"%v\"", browserBinary, browsers))
	}

	// Find browser type
	if browserType == "" {
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
					browserType = browser.Name

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
	if browserType == "" {
		crash("could not launch unknown browser type", fmt.Errorf("tried to launch prefered browser type (set with the HYDRAPP_TYPE env variable) \"%v\" and known types \"%v\"", browserType, browsers))
	}

	switch browserType {
	// Launch Chromium-like browser
	case browserTypeChromium:
		// Create a profile for the app
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			crash("could not get user's config directory", err)
		}
		userDataDir := filepath.Join(userConfigDir, id)

		// Create the browser instance
		execLine := append(
			browserBinary,
			append(
				[]string{
					"--name=" + name,
					"--class=" + name,
					"--user-data-dir=" + userDataDir,
					"--no-first-run",
					"--no-default-browser-check",
					"--app=" + url,
				},
				os.Args[1:]...,
			)...,
		)
		cmd := exec.Command(
			execLine[0],
			execLine[1:]...,
		)

		// Use system stdout, stderr and stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		// Start the browser
		if err := cmd.Run(); err != nil {
			crash("could not launch browser", err)
		}

		// Wait till lock for browser has been removed
		waitForLock(filepath.Join(userDataDir, "SingletonSocket"))

	// Launch Firefox-like browser
	case browserTypeFirefox:
		// Create a profile for the app
		execLine := append(
			browserBinary,
			[]string{
				"--createprofile",
				id,
			}...,
		)

		if output, err := exec.Command(
			execLine[0],
			execLine[1:]...,
		).CombinedOutput(); err != nil {
			crash("could not create profile", fmt.Errorf("%v: %v", err, string(output)))
		}

		// Get the user's home directory in which the profiles can be found
		home, err := os.UserHomeDir()
		if err != nil {
			crash("could not get user's home directory", err)
		}

		// Get the profile's directory
		firefoxDir := filepath.Join(home, ".mozilla", "firefox")
		if runtime.GOOS == osTypeWindows || runtime.GOOS == osTypeMacOS {
			userConfigDir, err := os.UserConfigDir()
			if err != nil {
				crash("could not get user's config directory", err)
			}

			if runtime.GOOS == osTypeWindows {
				firefoxDir = filepath.Join(userConfigDir, "Mozilla", "Firefox", "Profiles")
			} else {
				firefoxDir = filepath.Join(userConfigDir, "Firefox", "Profiles")
			}
		}

		filesInFirefoxDir, err := ioutil.ReadDir(firefoxDir)
		if err != nil {
			crash("could not get files in profiles directory", err)
		}

		profileSuffix := ""
		for _, file := range filesInFirefoxDir {
			if strings.HasSuffix(file.Name(), id) {
				profileSuffix = file.Name()

				break
			}
		}

		if profileSuffix == "" {
			crash("could not find profile directory generated by Firefox", fmt.Errorf("the profile's directory does not exist"))
		}

		profileDir := filepath.Join(firefoxDir, profileSuffix)
		if err := os.Setenv("PROFILE_DIR", profileDir); err != nil {
			crash("could not set profile directory", err)
		}

		// Add PWA styling using the preferences file
		prefsFile, err := os.OpenFile(filepath.Join(profileDir, "prefs.js"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
		if err != nil {
			crash("could not open preferences file", err)
		}
		defer prefsFile.Close()

		prefsFileContent, err := ioutil.ReadAll(prefsFile)
		if err != nil {
			crash("could not read preferences file", err)
		}

		for _, line := range strings.Split(firefoxPrefs, "\n") {
			if !strings.Contains(string(prefsFileContent), line) {
				if _, err := prefsFile.WriteString("\n" + line); err != nil {
					crash("could not write to preferences file", err)
				}
			}
		}

		// Add PWA styling using the userChrome file
		userChromeDir := filepath.Join(profileDir, "chrome")
		if err := os.MkdirAll(userChromeDir, 0755); err != nil {
			crash("could not create user chrome directory", err)
		}

		userChromeFile, err := os.OpenFile(filepath.Join(userChromeDir, "userChrome.css"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0664)
		if err != nil {
			crash("could not create user chrome file", err)
		}
		defer userChromeFile.Close()

		userChromeContent, err := ioutil.ReadAll(userChromeFile)
		if err != nil {
			crash("could not read user chrome file", err)
		}

		if !strings.Contains(string(userChromeContent), firefoxUserChrome) {
			if _, err := userChromeFile.WriteString("\n" + firefoxUserChrome); err != nil {
				crash("could not write to user chrome file", err)
			}
		}

		// Create the browser instance
		execLine = append(
			browserBinary,
			append(
				[]string{
					"--name=" + name,
					"--class=" + name,
					"--new-window",
					"--no-first-run",
					"-P",
					id,
					url,
				},
				os.Args[1:]...,
			)...,
		)

		cmd := exec.Command(
			execLine[0],
			execLine[1:]...,
		)

		// Use system stdout, stderr and stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		// Start the browser
		if err := cmd.Run(); err != nil {
			crash("could not launch browser", err)
		}

		// Wait till lock for browser has been removed
		waitForLock(filepath.Join(profileDir, "cookies.sqlite-wal"))

	// Launch Epiphany-like browser
	case browserTypeEpiphany:
		// Get the user's home directory in which the profiles should be created
		home, err := os.UserHomeDir()
		if err != nil {
			crash("could not get user's home directory", err)
		}

		// Create the profile directory
		epiphanyID := "org.gnome.Epiphany.WebApp-" + id
		profileDir := filepath.Join(home, ".local", "share", epiphanyID)

		if err := os.MkdirAll(filepath.Join(profileDir, ".app"), 0755); err != nil {
			crash("could not create profile directory", err)
		}

		// Create the .desktop file
		if err := ioutil.WriteFile(
			filepath.Join(profileDir, epiphanyID+".desktop"),
			[]byte(fmt.Sprintf(
				epiphanyDesktopFileTemplate,
				name,
				name,
				profileDir,
				url,
				name,
				epiphanyID,
			)),
			0664); err != nil {
			crash("could not write desktop file", err)
		}

		// Create the browser instance
		execLine := append(
			browserBinary,
			append(
				[]string{
					"--name=" + name,
					"--class=" + name,
					"--new-window",
					"--application-mode",
					"--profile=" + profileDir,
					url,
				},
				os.Args[1:]...,
			)...,
		)

		cmd := exec.Command(
			execLine[0],
			execLine[1:]...,
		)

		// Use system stdout, stderr and stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		// Start the browser
		if err := cmd.Run(); err != nil {
			crash("could not launch browser", err)
		}

		// Wait till lock for browser has been removed
		waitForLock(filepath.Join(profileDir, "ephy-history.db-wal"))
	}
}

func crash(msg string, err error) {
	// Create user-friendly error message
	body := fmt.Sprintf(`%v has encountered a fatal error and can't continue. The error message is:

%v

The following information might help you in fixing the problem:

%v`,
		name,
		capitalize(msg),
		capitalize(err.Error()),
	)

	// Show error message visually using a dialog
	if err := zenity.Error(
		body,
		zenity.Title("Fatal Error"),
		zenity.Width(320),
	); err != nil {
		log.Println("could not display fatal error dialog:", err)
	}

	// Log error message and exit with non-zero exit code
	log.Fatalln(body)
}

func capitalize(msg string) string {
	// Capitalize the first letter of the message if it is longer than two characters
	if len(msg) >= 2 {
		return string(unicode.ToUpper([]rune(msg)[0])) + msg[1:]
	}

	return msg
}

func waitForLock(path string) {
	if _, err := os.Stat(path); err == nil {
		// Wait until browser has exited
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			crash("could not start lockfile watcher", err)
		}
		defer watcher.Close()

		done := make(chan struct{})
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}

					// Stop the app
					if event.Op&fsnotify.Remove == fsnotify.Remove {
						done <- struct{}{}

						return
					}

				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}

					crash("could not continue watching lockfile", err)
				}
			}
		}()

		err = watcher.Add(path)
		if err != nil {
			crash("could not watch lockfile", err)
		}

		<-done
	}
}
