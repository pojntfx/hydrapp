#!/bin/bash

set -e

# Setup GPG
echo "${GPG_KEY_PASSWORD}" | base64 -d >'/tmp/gpg-pass'
mkdir -p "${HOME}/.gnupg"
cat >"${HOME}/.gnupg/gpg.conf" <<EOT
yes
passphrase-file /tmp/gpg-pass
pinentry-mode loopback
EOT

echo "${GPG_KEY_CONTENT}" | base64 -d >'/tmp/private.gpg'
gpg --import /tmp/private.gpg

gpg --output /tmp/public.asc --armor --export

# Install Flatpak dependencies
flatpak remote-add --if-not-exists 'flathub' 'https://flathub.org/repo/flathub.flatpakrepo'

# Build app
export HOST_ARCH="$(uname -m)"
for ARCH in ${ARCHITECTURES}; do
  # Install pre-build SDKs
  flatpak install -y --arch="${ARCH}" "flathub org.freedesktop.Platform//21.08" "org.freedesktop.Sdk//21.08" "org.freedesktop.Sdk.Extension.golang//21.08"

  # Build SDK and export to repo
  flatpak-builder -y --arch="${ARCH}" --gpg-sign="${GPG_KEY_ID}" --repo="/dst" --force-clean --user --install "build-dir" "org.freedesktop.Sdk.Extension.ImageMagick.yaml"

  # Build app and export to repo
  flatpak-builder -y --arch="${ARCH}" --gpg-sign="${GPG_KEY_ID}" --repo="/dst" --force-clean "build-dir" "${APP_ID}.yaml"

  # Export `.flatpak` to out dir
  flatpak --arch="${ARCH}" --gpg-sign="${GPG_KEY_ID}" build-bundle "/dst" "out/${APP_ID}.linux-${ARCH}.flatpak" "${APP_ID}"
done

echo "[Flatpak Repo]
Title=Hydrapp Flatpak repo
Url=${BASE_URL}
Homepage=${BASE_URL}
Description=Flatpaks for Hydrapp
GPGKey=$(base64 -w 0 /tmp/public.gpg)
" >/dst/hydrapp.flatpakrepo
