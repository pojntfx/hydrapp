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

	spawnCmd  = "flatpak-spawn"
	spawnHost = "--host"

	browserTypeChromium = "chromium"
	browserTypeFirefox  = "firefox"
	browserTypeEpiphany = "epiphany"

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

type BrowserType struct {
	Name     string
	Binaries []string
}

var chromiumLikeBrowsers = BrowserType{
	Name: browserTypeChromium,
	Binaries: []string{
		"google-chrome",
		"google-chrome-stable",
		"google-chrome-beta",
		"google-chrome-unstable",
		"brave-browser",
		"brave-browser-stable",
		"brave-browser-beta",
		"brave-browser-nightly",
		"microsoft-edge",
		"microsoft-edge-beta",
		"microsoft-edge-dev",
		"ungoogled-chromium",
		"chromium-browser",
		"chromium",
	},
}

var firefoxLikeBrowsers = BrowserType{
	Name: browserTypeFirefox,
	Binaries: []string{
		"firefox",
	},
}

var epiphanyLikeBrowsers = BrowserType{
	Name: browserTypeEpiphany,
	Binaries: []string{
		"epiphany",
	},
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

	// Check if we are in Flatpak
	runningInFlatpak := false
	if _, err := exec.LookPath(spawnCmd); err == nil {
		runningInFlatpak = true
	}

	// Find browser binary
	// Order matters; whatever is discovered first will be used
	knownBrowserTypes := []BrowserType{chromiumLikeBrowsers, firefoxLikeBrowsers, epiphanyLikeBrowsers}
	if browserBinary[0] == "" {
	i:
		for _, knownBrowserType := range knownBrowserTypes {
			for _, knownBrowserBinary := range knownBrowserType.Binaries {
				if runningInFlatpak {
					// Find supported browser from Flatpak
					if err := exec.Command(spawnCmd, spawnHost, "which", knownBrowserBinary).Run(); err == nil {
						browserBinary = []string{spawnCmd, spawnHost, knownBrowserBinary}

						break i
					}
				} else {
					// Find supported browser in native install
					if _, err := exec.LookPath(knownBrowserBinary); err == nil {
						browserBinary = []string{knownBrowserBinary}

						break i
					}
				}
			}
		}
	} else {
		if runningInFlatpak {
			// Allow same override in Flatpak by spawning the prefered browser on the host
			browserBinary = append([]string{spawnCmd, spawnHost}, browserBinary...)
		}
	}

	// Abort if browser binary could not be found
	if browserBinary[0] == "" {
		crash("could not find a supported browser", fmt.Errorf("tried to launch preferred browser binary (set with the HYDRAPP_BROWSER env variable) \"%v\" and known binaries \"%v\"", browserBinary, knownBrowserTypes))
	}

	// Find browser type
	if browserType == "" {
	j:
		for _, knownBrowserType := range knownBrowserTypes {
			for _, knownBrowserBinary := range knownBrowserType.Binaries {
				if browserBinary[0] == knownBrowserBinary {
					browserType = knownBrowserType.Name

					break j
				}
			}
		}
	}

	// Abort if browser type could not be found
	if browserType == "" {
		crash("could not launch unknown browser type", fmt.Errorf("tried to launch prefered browser type (set with the HYDRAPP_TYPE env variable) \"%v\" and known types \"%v\"", browserType, knownBrowserTypes))
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
		if err := exec.Command(browserBinary[0], "--createprofile", id).Run(); err != nil {
			crash("could not create profile", err)
		}

		// Get the user's home directory in which the profiles can be found
		home, err := os.UserHomeDir()
		if err != nil {
			crash("could not get user's home directory", err)
		}

		// Get the profile's directory
		firefoxDir := filepath.Join(home, ".mozilla", browserBinary[0])
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
		execLine := append(
			browserBinary,
			append(
				[]string{
					"--name=" + name,
					"--class=" + name,
					"--new-window",
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
