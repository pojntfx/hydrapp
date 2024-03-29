#!/bin/bash

set -e

# Setup PGP
echo "${PGP_KEY_PASSWORD}" | base64 -d >'/tmp/pgp-pass'
mkdir -p "${HOME}/.gnupg"
cat >"${HOME}/.gnupg/gpg.conf" <<EOT
yes
passphrase-file /tmp/pgp-pass
pinentry-mode loopback
EOT

echo "${PGP_KEY}" | base64 -d >'/tmp/private.pgp'
gpg --import /tmp/private.pgp

# Prepare build environment
export BASEDIR="${PWD}/${GOMAIN}"

gpg --output /tmp/repo.asc --armor --export

# Install Flatpak dependencies
flatpak remote-add --if-not-exists 'flathub' 'https://flathub.org/repo/flathub.flatpakrepo'

# See https://github.com/pojntfx/bagop/blob/main/main.go#L45
export DEBARCH="${GOARCH}"
if [ "${ARCHITECTURE}" = "386" ]; then
    export DEBARCH="i686"
elif [ "${ARCHITECTURE}" = "amd64" ]; then
    export DEBARCH="x86_64"
elif [ "${ARCHITECTURE}" = "arm" ]; then
    export DEBARCH="armv7l"
elif [ "${ARCHITECTURE}" = "arm64" ]; then
    export DEBARCH="aarch64"
fi

# Install pre-build SDKs
flatpak install -y --arch="${DEBARCH}" 'flathub' "org.freedesktop.Platform//23.08" "org.freedesktop.Sdk//23.08" "org.freedesktop.Sdk.Extension.golang//23.08" "org.freedesktop.Sdk.Extension.node18//23.08"

# Build SDK and export to repo
flatpak-builder -y --arch="${DEBARCH}" --gpg-sign="$(echo ${PGP_KEY_ID} | base64 -d)" --repo='/dst' --force-clean --user --install "build-dir" "${GOMAIN}/org.freedesktop.Sdk.Extension.ImageMagick.yaml"

# Build app and export to repo
flatpak-builder -y --arch="${DEBARCH}" --gpg-sign="$(echo ${PGP_KEY_ID} | base64 -d)" --repo='/dst' --force-clean "build-dir" "${GOMAIN}/${APP_ID}.yaml"

# Export `.flatpak` to out dir
flatpak --arch="${DEBARCH}" --gpg-sign="$(echo ${PGP_KEY_ID} | base64 -d)" build-bundle '/dst' "/dst/${APP_ID}.linux-${DEBARCH}.flatpak" "${APP_ID}"

echo "[Flatpak Repo]
Title=Hydrapp Flatpak repo
Url=${BASE_URL}
Homepage=${BASE_URL}
Description=Flatpaks for Hydrapp
GPGKey=$(base64 -w 0 /tmp/repo.asc)
" >/dst/hydrapp.flatpakrepo
