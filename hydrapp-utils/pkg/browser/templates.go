package browser

const (
	prefsJSContent = `user_pref("toolkit.legacyUserProfileCustomizations.stylesheets", true);
user_pref("browser.link.open_newwindow.restriction", 0);
user_pref("browser.tabs.firefox-view", false);
user_pref("browser.link.open_newwindow", 1);`
	userChromeCSSContent = `@namespace url("http://www.mozilla.org/keymaster/gatekeeper/there.is.only.xul");

#tabbrowser-tabs {
  visibility: collapse !important;
}

browser {
  margin-right: 0px;
  margin-bottom: 0px;
}

#main-window[chromehidden*="toolbar"] #nav-bar {
  visibility: collapse !important;
}

#nav-bar {
  margin-top: 0;
  margin-bottom: -42px;
  z-index: -100;
}

#PersonalToolbar {
  display: none;
}
`

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
