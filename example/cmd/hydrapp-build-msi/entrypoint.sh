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

# Create icons
convert -resize 'x16' -gravity 'center' -crop '16x16+0+0' -flatten -colors '256' -background 'transparent' 'icon.png' '/dst/icon.ico'

# Build app
export GOOS="windows"
for ARCH in ${ARCHITECTURES}; do
  export GOARCH="${ARCH}"

  # See https://github.com/pojntfx/bagop/blob/main/main.go#L45
  export DEBARCH="${GOARCH}"
  if [ "${ARCH}" = "386" ]; then
    export DEBARCH="i686"
  elif [ "${ARCH}" = "amd64" ]; then
    export DEBARCH="x86_64"
  elif [ "${ARCH}" = "arm" ]; then
    export DEBARCH="armv7l"
  elif [ "${ARCH}" = "arm64" ]; then
    export DEBARCH="aarch64"
  fi

  GOFLAGS='-tags=selfupdate' go build -o "/dst/${APP_ID}.${GOOS}-${DEBARCH}.exe" .
  gpg --detach-sign --armor "/dst/${APP_ID}.${GOOS}-${DEBARCH}.exe"

  wixl -o "/dst/${APP_ID}.${GOOS}-${DEBARCH}.msi" <(sed "s@Source=\"${APP_ID}.exe\"@Source=\"/dst/${APP_ID}.${GOOS}-${DEBARCH}.exe\"@g" ${APP_ID}.wxl | sed 's@SourceFile="icon.ico"@SourceFile="/dst/icon.ico"@g' | sed 's@Source="icon.ico"@Source="/dst/icon.ico"@g')
  gpg --detach-sign --armor "/dst/${APP_ID}.${GOOS}-${DEBARCH}.msi"
done

rm '/dst/icon.ico'
