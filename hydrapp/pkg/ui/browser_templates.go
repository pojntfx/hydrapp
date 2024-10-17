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
user_pref("browser.link.open_newwindow", 1);`
	userChromeCSSContent = `@namespace url("http://www.mozilla.org/keymaster/gatekeeper/there.is.only.xul");

#nav-bar,
#urlbar-container,
#searchbar,
#PersonalToolbar,
#tabbrowser-tabs,
#TabsToolbar #firefox-view-button,
#alltabs-button {
  visibility: collapse !important;
}

#navigator-toolbox {
    border: 0 !important;
}

#navigator-toolbox[inFullscreen] {
    position: relative;
    z-index: 1;
    height: 3px;
    margin-bottom: -3px;
    opacity: 0;
    overflow: hidden;
}

#navigator-toolbox[inFullscreen]:hover {
    height: auto;
    margin-bottom: 0px;
    opacity: 1;
    overflow: show;
}

#content-deck[inFullscreen]{
    position:relative;
    z-index: 0;
}
`

	epiphanyDesktopFileTemplate = `[Desktop Entry]
Name=%v
Exec=epiphany --new-window --application-mode --profile="%v" "%v"
StartupNotify=true
Terminal=false
Type=Application
Categories=GNOME;GTK;
StartupWMClass=%v
X-Purism-FormFactor=Workstation;Mobile;`
)
