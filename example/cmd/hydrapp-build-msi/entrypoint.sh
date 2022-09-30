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

# Install MSYS2 packages
if [ "${MSYS2PACKAGES}" != "" ]; then
  wine64 bash.exe -c "pacman --noconfirm --ignore pacman --needed -S ${MSYS2PACKAGES}"
fi

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

  if [ "${ARCH}" = "amd64" ]; then
    for UNIXPATH in $(find . -type f); do
      export WINDOWSPATH="$(sed s,/,\\\\,g <<<$(echo ${UNIXPATH} | sed s,^./,,g))"
      export UUID="$(uuid)"

      echo "<File Id=\"${UUID}\" Source=\"${WINDOWSPATH}\" KeyPath=\"yes\" DiskId=\"1\" />"

      if [[ "${UNIXPATH}" == *${APP_ID} ]]; then
        echo "<Shortcut Id=\"${UUID}\" Advertise=\"yes\" Icon=\"icon.ico\" Name=\"${APP_NAME}\" Directory=\"ProgramMenuFolder\" WorkingDirectory=\"INSTALLDIR\" Description=\"${APP_NAME}\" />"
        echo "<Shortcut Id=\"${UUID}\" Advertise=\"yes\" Icon=\"icon.ico\" Name=\"${APP_NAME}\" Directory=\"DesktopFolder\" WorkingDirectory=\"INSTALLDIR\" Description=\"${APP_NAME}\" />"
      fi
    done
  else
    GOFLAGS='-tags=selfupdate' go build -o "/dst/${APP_ID}.${GOOS}-${DEBARCH}.exe" .
    gpg --detach-sign --armor "/dst/${APP_ID}.${GOOS}-${DEBARCH}.exe"

    wixl -o "/dst/${APP_ID}.${GOOS}-${DEBARCH}.msi" <(sed "s@Source=\"${APP_ID}.exe\"@Source=\"/dst/${APP_ID}.${GOOS}-${DEBARCH}.exe\"@g" ${APP_ID}.wxl | sed 's@SourceFile="icon.ico"@SourceFile="/dst/icon.ico"@g' | sed 's@Source="icon.ico"@Source="/dst/icon.ico"@g')
    gpg --detach-sign --armor "/dst/${APP_ID}.${GOOS}-${DEBARCH}.msi"
  fi
done

rm '/dst/icon.ico'
