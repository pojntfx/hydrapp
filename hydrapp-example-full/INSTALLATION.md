# Hydrapp Example Full (Main) Installation Instructions

## Android

1. Download F-Droid from [f-droid.org](https://f-droid.org/)
2. Open F-Droid and go to `Settings`, then `Repositories`
3. Add a new repository with the address `https://pojntfx.github.io/hydrapp/hydrapp-example-full/apk/main/`
4. Search for `Hydrapp Example Full (Main)` in F-Droid and tap on `Install`

Hydrapp Example Full (Main) should now be installed and receive updates automatically.

## Debian Sid on `amd64`

To install the prebuilt binary package, run the following:

```shell
pkexec sudo bash - <<'EOT'
mkdir -p /usr/local/share/keyrings
curl -Lo https://pojntfx.github.io/hydrapp/hydrapp-example-full/deb/debian/sid/x86_64/main/repo.asc | gpg --dearmor --output /usr/local/share/keyrings/hydrapp-example-full-main.gpg -

cat >/etc/apt/sources.list.d/hydrapp-example-full-main.list <<EOA
deb [signed-by=/usr/local/share/keyrings/hydrapp-example-full-main.gpg] https://pojntfx.github.io/hydrapp/hydrapp-example-full/deb/debian/sid/x86_64/main/ sid main
deb-src [signed-by=/usr/local/share/keyrings/hydrapp-example-full-main.gpg] https://pojntfx.github.io/hydrapp/hydrapp-example-full/deb/debian/sid/x86_64/main/ sid main
EOA

apt update

apt install -y 'com.pojtinger.felicitas.hydrapp.example.full.main'
EOT
```

To install the source package, build the binary package locally and install it, run the following:

```shell
pkexec sudo bash - <<'EOT'
mkdir -p /usr/local/share/keyrings
curl -Lo https://pojntfx.github.io/hydrapp/hydrapp-example-full/deb/debian/sid/x86_64/main/repo.asc | gpg --dearmor --output /usr/local/share/keyrings/hydrapp-example-full-main.gpg -

cat >/etc/apt/sources.list.d/hydrapp-example-full-main.list <<EOA
deb [signed-by=/usr/local/share/keyrings/hydrapp-example-full-main.gpg] https://pojntfx.github.io/hydrapp/hydrapp-example-full/deb/debian/sid/x86_64/main/ sid main
deb-src [signed-by=/usr/local/share/keyrings/hydrapp-example-full-main.gpg] https://pojntfx.github.io/hydrapp/hydrapp-example-full/deb/debian/sid/x86_64/main/ sid main
EOA

apt update

apt install -y dpkg-dev
apt build-dep -y 'com.pojtinger.felicitas.hydrapp.example.full.main'
apt source -y --build 'com.pojtinger.felicitas.hydrapp.example.full.main'
apt install -y ./com.pojtinger.felicitas.hydrapp.example.full.main_*.deb
EOT
```

Hydrapp Example Full (Main) should now be installed and receive updates automatically.

## macOS

1. Download the `.app` from [https://pojntfx.github.io/hydrapp/hydrapp-example-full/dmg/main/com.pojtinger.felicitas.hydrapp.example.full.main.darwin.dmg](https://pojntfx.github.io/hydrapp/hydrapp-example-full/dmg/main/com.pojtinger.felicitas.hydrapp.example.full.main.darwin.dmg)
2. Mount `com.pojtinger.felicitas.hydrapp.example.full.main.darwin.dmg` by opening it
3. Drag `Hydrapp Example Full (Main)` to the `Applications` folder

Hydrapp Example Full (Main) should now be installed and receive updates automatically.

## Linux Universal (Flatpak) on `amd64`

To install the package, run the following:

```shell
flatpak remote-add hydrapp-example-full-main --from 'https://pojntfx.github.io/hydrapp/hydrapp-example-full/flatpak/x86_64/main/hydrapp.flatpakrepo'
flatpak install -y 'com.pojtinger.felicitas.hydrapp.example.full.main'
```

Hydrapp Example Full (Main) should now be installed and receive updates automatically.

## Windows on `amd64`

1. Download the installer from [https://pojntfx.github.io/hydrapp/hydrapp-example-full/msi/x86_64/main/com.pojtinger.felicitas.hydrapp.example.full.main.windows-x86_64.msi](https://pojntfx.github.io/hydrapp/hydrapp-example-full/msi/x86_64/main/com.pojtinger.felicitas.hydrapp.example.full.main.windows-x86_64.msi)
2. Install the application by opening the downloaded file

Hydrapp Example Full (Main) should now be installed and receive updates automatically.

## Fedora 36 on `amd64`

To install the prebuilt binary package, run the following:

```shell
pkexec sudo bash - <<'EOT'
dnf config-manager --add-repo 'https://pojntfx.github.io/hydrapp/hydrapp-example-full/rpm/fedora/36/x86_64/main/repodata/hydrapp.repo'
dnf install -y 'com.pojtinger.felicitas.hydrapp.example.full.main'
EOT
```

To install the source package, build the binary package locally and install it, run the following:

```shell
pkexec sudo bash - <<'EOT'
dnf config-manager --add-repo 'https://pojntfx.github.io/hydrapp/hydrapp-example-full/rpm/fedora/36/x86_64/main/repodata/hydrapp.repo'

dnf install -y rpm-build
dnf download --source -y 'com.pojtinger.felicitas.hydrapp.example.full.main'
dnf builddep -y com.pojtinger.felicitas.hydrapp.example.full.main-*.rpm
rpmbuild --rebuild com.pojtinger.felicitas.hydrapp.example.full.main-*.rpm
dnf install -y ~/rpmbuild/RPMS/"$(uname -m)"/com.pojtinger.felicitas.hydrapp.example.full.main-*.rpm
EOT
```

Hydrapp Example Full (Main) should now be installed and receive updates automatically.
