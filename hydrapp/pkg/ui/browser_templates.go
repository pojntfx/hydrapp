package ui

const (
	prefsJSContent = `user_pref("toolkit.legacyUserProfileCustomizations.stylesheets", true);
user_pref("browser.link.open_newwindow.restriction", 0);
user_pref("browser.tabs.firefox-view", false);
user_pref("datareporting.policy.firstRunURL", "");
user_pref("browser.tabs.tabmanager.enabled", false);
user_pref("browser.shell.checkDefaultBrowser", false);
user_pref("datareporting.policy.dataSubmissionPolicyBypassNotification", true);
user_pref("browser.tabs.warnOnClose", false);
user_pref("browser.tabs.warnOnCloseOtherTabs", false);
user_pref("browser.warnOnQuit", false);
user_pref("browser.warnOnQuitShortcut", false);
user_pref("browser.newtabpage.activity-stream.asrouter.userprefs.cfr.features", false);
user_pref("disableResetPrompt", true);
user_pref("trailhead.firstrun.branches", "nofirstrun-empty");
user_pref("browser.aboutwelcome.enabled", false);
user_pref("browser.link.open_newwindow", 1);
user_pref("full-screen-api.macos-native-full-screen", true);
user_pref("browser.fullscreen.autohide", true);
user_pref("browser.sessionstore.resume_from_crash", false);
user_pref("ui.key.menuAccessKeyFocuses", false);`
	userChromeCSSContent = `@namespace url("http://www.mozilla.org/keymaster/gatekeeper/there.is.only.xul");

#TabsToolbar-customization-target {
  visibility: hidden !important;
}

#PersonalToolbar,
#tabbrowser-tabs,
#urlbar {
  display: none !important;
}

#nav-bar {
  z-index: -1;
  height: 0;
  width: 0;
  min-height: 0 !important;
  min-width: 0 !important;
  border: 0 !important;
  overflow: hidden;
}`

	epiphanyDesktopFileTemplate = `[Desktop Entry]
Exec=%v --new-window --application-mode --profile="%v" "%v"
StartupNotify=true
Terminal=false
Type=Application
Categories=GNOME;GTK;
StartupWMClass=%v
X-Purism-FormFactor=Workstation;Mobile;
Name=%v
Icon=%v
NoDisplay=true`
)
