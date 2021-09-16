# Desktop Support

Some examples on how to launch browsers in a "single-app"/PWA mode; eventually, a Go-based library for this functionality could be created.

## Preparation

```shell
$ export HYDRAP_URL=https://go-app.dev/
$ export HYDRAP_NAME="Go-App"
$ export HYDRAP_ID=go-app
```

## Chrome/Chromium/Edge/Brave etc.

```shell
$ google-chrome-stable --no-first-run --no-default-browser-check --user-data-dir="/tmp/${HYDRAP_ID}" --app="${HYDRAP_URL}" --class="${HYDRAP_NAME}"
```

## Firefox

```shell
$ firefox --createprofile "${HYDRAP_ID}"
$ export PROFILE_DIR=~/.mozilla/firefox/$(ls ~/.mozilla/firefox/ | grep ${HYDRAP_ID})
$ printf 'user_pref("toolkit.legacyUserProfileCustomizations.stylesheets", true);\nuser_pref("browser.tabs.drawInTitlebar", false);' >> ${PROFILE_DIR}/prefs.js
$ mkdir -p ${PROFILE_DIR}/chrome
$ cat <<EOT> ${PROFILE_DIR}/chrome/userChrome.css
TabsToolbar {
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
}
EOT
$ firefox -P "${HYDRAP_ID}" "${HYDRAP_URL}" --class="${HYDRAP_NAME}"
```

## GNOME Web

```shell
# For Flatpak
$ alias epiphany="flatpak run org.gnome.Epiphany"
$ export PROFILE_DIR=$HOME/Downloads/org.gnome.Epiphany.WebApp-${HYDRAP_ID}
# For native:
$ export PROFILE_DIR=$HOME/.local/share/org.gnome.Epiphany.WebApp-${HYDRAP_ID}
# For both:
$ mkdir -p ${PROFILE_DIR}/.app
$ cat <<EOT>${PROFILE_DIR}/org.gnome.Epiphany.WebApp-${HYDRAP_ID}.desktop
[Desktop Entry]
Name=${HYDRAP_NAME}
Exec=epiphany --application-mode --profile="${PROFILE_DIR}" "${HYDRAP_URL}"
StartupNotify=true
Terminal=false
Type=Application
Categories=GNOME;GTK;
StartupWMClass=org.gnome.Epiphany.WebApp-${HYDRAP_ID}
X-Purism-FormFactor=Workstation;Mobile;
EOT
$ epiphany --application-mode --profile=${PROFILE_DIR} "${HYDRAP_URL}" --class="${HYDRAP_NAME}"
```

## Example App Installation on Debian and Ubuntu

```shell
apt update
apt install -y ca-certificates gnupg2 lsb-release

gpg --keyserver keyserver.ubuntu.com --recv-keys 638840CAE7660B1B69ADEE9041DDCDD3AFF03AC7

mkdir -p /usr/local/share/keyrings
gpg --output /usr/local/share/keyrings/hydrapp.gpg --export 638840CAE7660B1B69ADEE9041DDCDD3AFF03AC7

cat >/etc/apt/sources.list.d/hydrapp.list <<EOT
deb [signed-by=/usr/local/share/keyrings/hydrapp.gpg] https://pojntfx.github.io/hydrapp/apt/$(lsb_release -i -s | tr '[:upper:]' '[:lower:]')/ $(lsb_release -c -s) main
deb-src [signed-by=/usr/local/share/keyrings/hydrapp.gpg] https://pojntfx.github.io/hydrapp/apt/$(lsb_release -i -s | tr '[:upper:]' '[:lower:]')/ $(lsb_release -c -s) main
EOT

apt update

# Install binary package
apt install -y com.pojtinger.felicitas.hydrapp.example

# Alternatively install from source
apt install -y dpkg-dev

apt -y build-dep com.pojtinger.felicitas.hydrapp.example
apt source -y --build com.pojtinger.felicitas.hydrapp.example
apt install -y ./com.pojtinger.felicitas.hydrapp.example_*.deb
```

## Example App Installation on Fedora

```shell
dnf install -y gnupg2 'dnf-command(config-manager)'

gpg --keyserver keyserver.ubuntu.com --recv-keys 638840CAE7660B1B69ADEE9041DDCDD3AFF03AC7

mkdir -p /usr/local/share/keyrings
gpg --output /usr/local/share/keyrings/hydrapp.gpg --export 638840CAE7660B1B69ADEE9041DDCDD3AFF03AC7

echo "[hydrapp-repo]
name=Hydrapp YUM repo
baseurl=https://pojntfx.github.io/hydrapp/yum/fedora-\$releasever
enabled=1
gpgcheck=1
gpgkey=file:///usr/local/share/keyrings/hydrapp.gpg" >/tmp/hydrapp.repo
dnf config-manager --add-repo /tmp/hydrapp.repo

# Install binary package
dnf install -y com.pojtinger.felicitas.hydrapp.example

# Alternatively install from source
sudo dnf install -y rpm-build

dnf download --source -y com.pojtinger.felicitas.hydrapp.example
dnf builddep -y com.pojtinger.felicitas.hydrapp.example-*.rpm
rpmbuild --rebuild com.pojtinger.felicitas.hydrapp.example-*.rpm
sudo dnf install -y ~/rpmbuild/RPMS/"$(uname -m)"/com.pojtinger.felicitas.hydrapp.example-*.rpm
```
