{{ $appID := .AppID }}
{{ $appName := .AppName }}

# {{ $appName }} Installation Instructions

## Android

1. Download F-Droid from [f-droid.org](https://f-droid.org/)
2. Open F-Droid and go to `Settings`, then `Repositories`
3. Add a new repository with the address `{{ .AndroidRepoURL }}`
4. Search for `{{ $appName }}` in F-Droid and tap on `Install`

{{ $appName }} should now be installed and receive updates automatically.

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

{{ $appName }} should now be installed and receive updates automatically.

## macOS

1. Download the `.app` from [{{ .MacOSBinaryURL }}]({{ .MacOSBinaryURL }})
2. Mount `{{ .MacOSBinaryName }}` by opening it
3. Drag `{{ $appName }}` to the `Applications` folder

{{ $appName }} should now be installed and receive updates automatically.

{{ range .Flatpaks }}

## Linux Universal (Flatpak) on `{{ .Architecture }}`

To install the package, run the following:

```shell
flatpak remote-add '{{ $appID }}' --from '{{ .URL }}'
flatpak install -y '{{ $appID }}'
```

{{ $appName }} should now be installed and receive updates automatically.
{{ end }}

{{ range .MSIs }}

## Windows on `{{ .Architecture }}`

1. Download the installer from [{{ .URL }}]({{ .URL }})
2. Install the application by opening the downloaded file

{{ $appName }} should now be installed and receive updates automatically.

{{ end }}

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

{{ $appName }} should now be installed and receive updates automatically.

## Other Platforms

{{ $appName }} is also available as a single static binary. To install it, download the binary for your operating system and processor architecture from [https://pojntfx.github.io/hydrapp/hydrapp-example-full/binaries/main/](https://pojntfx.github.io/hydrapp/hydrapp-example-full/binaries/main/). They include a self-update mechanism, so you should be receiving updates automatically.
