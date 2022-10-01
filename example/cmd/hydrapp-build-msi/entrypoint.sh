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
mkdir -p '/tmp/out'
convert -resize 'x16' -gravity 'center' -crop '16x16+0+0' -flatten -colors '256' -background 'transparent' 'icon.png' '/tmp/out/icon.ico'

# Build app
export GOOS="windows"
export BASEDIR="${PWD}"
for ARCH in ${ARCHITECTURES}; do
  cd "${BASEDIR}"

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
    GOPATH='/root/.wine/drive_c/go' go mod download -x

    cp -r . '/root/.wine/drive_c/users/root/Documents/go-workspace'
    rm -rf '/root/.wine/drive_c/users/root/Documents/go-workspace/out'

    wine64 bash.exe -c "export PATH=$PATH:/mingw64/bin GOPATH=/c/go GOROOT=/mingw64/lib/go TMP=/c/tmp TEMP=/c/tmp GOARCH=amd64 CGO_ENABLED=1 && cd /c/users/root/Documents/go-workspace && go build -buildvcs=false -ldflags='-linkmode=external' -x -v -o out/${APP_ID}.${GOOS}-${DEBARCH}.exe ."

    # Copy binaries to staging directory
    yes | cp -rf /root/.wine/drive_c/users/root/Documents/go-workspace/out/* '/tmp/out'

    cp -r /root/.wine/drive_c/msys64/mingw64/* '/tmp/out'
  else
    GOFLAGS='-tags=selfupdate' go build -o "/tmp/out/${APP_ID}.${GOOS}-${DEBARCH}.exe" .
  fi

  cd '/tmp/out'

  # Prepare WiX installer
  export STARTID=""
  for UNIXPATH in $(find . -type f); do
    export UUID="$(uuid)"

    echo "<File Id=\"${UUID}\" Name=\"${UNIXPATH}\" Source=\"${UNIXPATH}\" DiskId=\"1\" />\n" >>'/tmp/files'

    if [[ "${UNIXPATH}" == *${APP_ID}.${GOOS}-${DEBARCH}.exe ]]; then
      export STARTID="${UUID}"
    fi
  done

  wixl -v -o "/dst/${APP_ID}.${GOOS}-${DEBARCH}.msi" <(cat "${BASEDIR}/${APP_ID}.wxl" | perl -p -e 'use File::Slurp; my $text = read_file("/tmp/files"); s+<HydrappFiles />+$text+g' | perl -p -e 's+{{ StartID }}+$ENV{STARTID}+g')
  gpg --detach-sign --armor "/dst/${APP_ID}.${GOOS}-${DEBARCH}.msi"
done
