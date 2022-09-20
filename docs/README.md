# Installation

## Example App Installation on Debian and Ubuntu

```shell
apt update
apt install -y ca-certificates lsb-release

mkdir -p /usr/local/share/keyrings
curl -L -o /tmp/hydrapp.asc https://pojntfx.github.io/hydrapp/apt/repo.asc
gpg --dearmor --output /usr/local/share/keyrings/hydrapp.gpg /tmp/hydrapp.asc

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

## Example App Installation on Fedora, CentOS and OpenSUSE

```shell
source /etc/os-release

if [ "${ID}" = "opensuse-tumbleweed" ]; then
    zypper install -y dirmngr
else
    dnf install -y 'dnf-command(config-manager)'
fi

if [ "${ID}" = "opensuse-tumbleweed" ]; then
    zypper addrepo https://pojntfx.github.io/hydrapp/yum/opensuse-tumbleweed/repodata/hydrapp.repo
else
    dnf config-manager --add-repo https://pojntfx.github.io/hydrapp/yum/$([ "${ID}" = "centos" ] && echo "epel" || echo "${ID}")-"${VERSION_ID}"/repodata/hydrapp.repo
fi

# Install binary package
if [ "${ID}" = "opensuse-tumbleweed" ]; then
    zypper install -y com.pojtinger.felicitas.hydrapp.example
else
    dnf install -y com.pojtinger.felicitas.hydrapp.example
fi

# Alternatively install from source
if [ "${ID}" = "opensuse-tumbleweed" ]; then
    zypper install -y rpm-build

    zypper -n source-install com.pojtinger.felicitas.hydrapp.example

    rpmbuild -ba /usr/src/packages/SPECS/com.pojtinger.felicitas.hydrapp.example.spec

    zypper --no-gpg-checks install -y /usr/src/packages/RPMS/"$(uname -m)"/com.pojtinger.felicitas.hydrapp.example-*.rpm
else
    dnf install -y rpm-build

    dnf download --source -y com.pojtinger.felicitas.hydrapp.example
    dnf builddep -y com.pojtinger.felicitas.hydrapp.example-*.rpm
    rpmbuild --rebuild com.pojtinger.felicitas.hydrapp.example-*.rpm
    dnf install -y ~/rpmbuild/RPMS/"$(uname -m)"/com.pojtinger.felicitas.hydrapp.example-*.rpm
fi
```

## Example App Installation with Flatpak

```shell
flatpak remote-add hydrapp --from https://pojntfx.github.io/hydrapp/flatpak/hydrapp.flatpakrepo

flatpak install -y com.pojtinger.felicitas.hydrapp.example
```

## Example App Installation for Android with F-Droid

1. Add `https://github.com/pojntfx/hydrapp/tree/gh-pages/fdroid` as a repo in F-Droid
2. Refresh indexes
3. Search and install "Hydrapp Example"

## Other Platforms

On other platforms, download the release from GitHub releases directly - it has self-update functionality included.
