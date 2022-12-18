{{ $appID := .AppID }}
{{ $appName := .AppName }}

# {{ $appName }} Installation Instructions

{{ if .RenderAPK }}
## Android

1. Download F-Droid from [f-droid.org](https://f-droid.org/)
2. Open F-Droid and go to `Settings`, then `Repositories`
3. Add a new repository with the address `{{ .AndroidRepoURL }}`
4. Search for `{{ $appName }}` in F-Droid and tap on `Install`

{{ $appName }} should now be installed and receive updates automatically.
{{ end }}

{{ range .DEBs }}

## {{ Titlecase .DistroName }} {{ Titlecase .DistroVersion }} on `{{ .Architecture }}`

To install the prebuilt binary package, run the following:

```shell
pkexec sudo bash - <<'EOT'
mkdir -p /usr/local/share/keyrings
curl -Lo {{ .URL }}/repo.asc | gpg --dearmor --output /usr/local/share/keyrings/{{ $appID }}.gpg -

cat >/etc/apt/sources.list.d/{{ $appID }}.list <<EOA
deb [signed-by=/usr/local/share/keyrings/{{ $appID }}.gpg] {{ .URL }} {{ .DistroName }} main
deb-src [signed-by=/usr/local/share/keyrings/{{ $appID }}.gpg] {{ .URL }} {{ .DistroName }} main
EOA

apt update

apt install -y '{{ $appID }}'
EOT
```

To install the source package, build the binary package locally and install it, run the following:

```shell
pkexec sudo bash - <<'EOT'
mkdir -p /usr/local/share/keyrings
curl -Lo {{ .URL }}repo.asc | gpg --dearmor --output /usr/local/share/keyrings/{{ $appID }}.gpg -

cat >/etc/apt/sources.list.d/{{ $appID }}.list <<EOA
deb [signed-by=/usr/local/share/keyrings/{{ $appID }}.gpg] {{ .URL }} {{ .DistroName }} main
deb-src [signed-by=/usr/local/share/keyrings/{{ $appID }}.gpg] {{ .URL }} {{ .DistroName }} main
EOA

apt update

apt install -y dpkg-dev
apt build-dep -y '{{ $appID }}'
apt source -y --build '{{ $appID }}'
apt install -y ./{{ $appID }}_*.deb
EOT
```

{{ $appName }} should now be installed and receive updates automatically.

{{ end }}

{{ range .RPMs }}

## {{ Titlecase .DistroName }} {{ Titlecase .DistroVersion }} on `{{ .Architecture }}`

To install the prebuilt binary package, run the following:

```shell
pkexec sudo bash - <<'EOT'
dnf config-manager --add-repo '{{ .URL }}'
dnf install -y '{{ $appID }}'
EOT
```

To install the source package, build the binary package locally and install it, run the following:

```shell
pkexec sudo bash - <<'EOT'
dnf config-manager --add-repo '{{ .URL }}'

dnf install -y rpm-build
dnf download --source -y '{{ $appID }}'
dnf builddep -y {{ $appID }}-*.rpm
rpmbuild --rebuild {{ $appID }}-*.rpm
dnf install -y ~/rpmbuild/RPMS/"$(uname -m)"/{{ $appID }}-*.rpm
EOT
```

{{ $appName }} should now be installed and receive updates automatically.

{{ end }}

{{ if .RenderDMG }}
## macOS

1. Download the `.app` from [{{ .MacOSBinaryURL }}]({{ .MacOSBinaryURL }})
2. Mount `{{ .MacOSBinaryName }}` by opening it
3. Drag `{{ $appName }}` to the `Applications` folder

{{ $appName }} should now be installed and receive updates automatically.
{{ end }}

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

{{ if .RenderBinaries }}
## Other Platforms

{{ $appName }} is also available as a single static binary. To install it, download the binary for your operating system and processor architecture from [{{ .BinariesURL }}]({{ .BinariesURL }}). They include a self-update mechanism, so you should be receiving updates automatically.
{{ end }}