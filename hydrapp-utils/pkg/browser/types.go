package browser

const (
	browserTypeChromium = "chromium"
	browserTypeFirefox  = "firefox"
	browserTypeEpiphany = "epiphany"
	browserTypeLynx     = "lynx"

	browserTypeDummy = "dummy"
)

var ChromiumLikeBrowsers = Browser{
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
		{"com.google.Chrome"},
		{"com.google.ChromeDev"},
		{"com.brave.Browser"},
		{"com.microsoft.Edge"},
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

var FirefoxLikeBrowsers = Browser{
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

var EpiphanyLikeBrowsers = Browser{
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

var LynxLikeBrowsers = Browser{
	Name: browserTypeLynx,
	LinuxBinaries: [][]string{
		{"lynx"},
	},
	Flatpaks:        [][]string{},
	WindowsBinaries: [][]string{},
	MacOSBinaries:   [][]string{},
}
