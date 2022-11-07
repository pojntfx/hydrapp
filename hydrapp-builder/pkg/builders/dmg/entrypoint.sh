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

echo "${PGP_KEY_CONTENT}" | base64 -d >'/tmp/private.pgp'
gpg --import /tmp/private.pgp

# Prepare build environment
export BASEDIR="${PWD}/${GOMAIN}"

# Install MacPorts packages
if [ "${MACPORTS}" != "" ]; then
  osxcross-macports install --static ${MACPORTS}
fi

# Generate dependencies
GOFLAGS="${GOFLAGS}" sh -c "${GOGENERATE}"

# Create icons
mkdir -p '/tmp/out'
convert "${BASEDIR}/icon.png" '/tmp/out/icon.icns'

# Build app
export COMMIT_TIME_RFC3339="$(git log -1 --date=format:'%Y-%m-%dT%H:%M:%SZ' --format=%cd)"
export BRANCH_ID="stable"
if [ "$(git tag --points-at HEAD)" = "" ]; then
  export BRANCH_ID="$(git symbolic-ref --short HEAD)"
fi

export GOOS="darwin"
export BINARIES=""
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

  export CC="${DEBARCH}-apple-darwin20.4-cc"
  export CXX="${DEBARCH}-apple-darwin20.4-c++"
  export PKGCONFIG="${DEBARCH}-apple-darwin20.4-pkg-config"
  export PKG_CONFIG="${DEBARCH}-apple-darwin20.4-pkg-config"

  go build -ldflags="-X github.com/pojntfx/hydrapp/hydrapp-utils/pkg/update.CommitTimeRFC3339=${COMMIT_TIME_RFC3339} -X github.com/pojntfx/hydrapp/hydrapp-utils/pkg/update.BranchID=${BRANCH_ID}" -o "/tmp/${APP_ID}.${GOOS}-${DEBARCH}" "${GOMAIN}"
  rcodesign sign "/tmp/${APP_ID}.${GOOS}-${DEBARCH}"

  export BINARIES="${BINARIES} /tmp/${APP_ID}.${GOOS}-${DEBARCH}"
done

lipo -create -output "/tmp/out/${APP_ID}.${GOOS}" ${BINARIES}
gpg --detach-sign --armor "/tmp/out/${APP_ID}.${GOOS}"
rcodesign sign "/tmp/out/${APP_ID}.${GOOS}"

cp "${BASEDIR}/Info.plist" '/tmp/out/'

mkdir -p "/tmp/out/${APP_ID}.${GOOS}.dmg.mnt/${APP_NAME}.app/Contents/"{MacOS,Resources}
cp "/tmp/out/${APP_ID}.${GOOS}" "/tmp/out/${APP_ID}.${GOOS}.dmg.mnt/${APP_NAME}.app/Contents/MacOS/${APP_ID}"
cp "/tmp/out/${APP_ID}.${GOOS}.asc" "/tmp/out/${APP_ID}.${GOOS}.dmg.mnt/${APP_NAME}.app/Contents/MacOS/${APP_ID}.asc"
cp '/tmp/out/Info.plist' "/tmp/out/${APP_ID}.${GOOS}.dmg.mnt/${APP_NAME}.app/Contents"
cp '/tmp/out/icon.icns' "/tmp/out/${APP_ID}.${GOOS}.dmg.mnt/${APP_NAME}.app/Contents/Resources"

genisoimage -V "Install $(echo ${APP_NAME} | cut -c -24)" -D -R -apple -no-pad -o "/tmp/out/${APP_ID}.${GOOS}.dmg" "/tmp/out/${APP_ID}.${GOOS}.dmg.mnt"
gpg --detach-sign --armor "/tmp/out/${APP_ID}.${GOOS}.dmg"

cp "/tmp/out/${APP_ID}.${GOOS}.dmg" "/tmp/out/${APP_ID}.${GOOS}.dmg.asc" "/dst"

cd /dst

tree -J . -I 'index.html|index.json' | jq '.[0].contents' | jq ". |= map( . + {time: \"${COMMIT_TIME_RFC3339}\"} )" | tee 'index.json'
