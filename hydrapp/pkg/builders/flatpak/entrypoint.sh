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
flatpak install -y --arch="${DEBARCH}" 'flathub' "org.freedesktop.Platform//24.08" "org.freedesktop.Sdk//24.08" "org.freedesktop.Sdk.Extension.golang//24.08" "org.freedesktop.Sdk.Extension.node20//24.08"

# Install extra references
if [ "${FLATPAKREFS}" != "" ]; then
    flatpak install -y --arch="${DEBARCH}" 'flathub' ${FLATPAKREFS}
fi

# Build app and export to repo
flatpak-builder -y --arch="${DEBARCH}" --gpg-sign="$(echo ${PGP_KEY_ID} | base64 -d)" --repo='/hydrapp/dst' --force-clean "build-dir" "${GOMAIN}/${APP_ID}.json"

cp -f ${BASEDIR}/hydrapp.flatpakrepo /hydrapp/dst/hydrapp.flatpakrepo
perl -pi -e 'my $key = `base64 -w 0 /tmp/repo.asc`; chomp($key); s+\{ PGPPublicKey \}+$key+g' /hydrapp/dst/hydrapp.flatpakrepo

if [ "${DST_UID}" != "" ] && [ "${DST_GID}" != "" ]; then
    chown -R "${DST_UID}:${DST_GID}" /hydrapp/dst
fi
