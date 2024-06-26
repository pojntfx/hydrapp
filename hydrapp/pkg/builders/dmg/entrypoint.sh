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

# Configure Go
export GOPROXY='https://proxy.golang.org,direct'

# Generate dependencies
GOFLAGS="${GOFLAGS}" sh -c "${GOGENERATE}"

# Create icons
mkdir -p '/tmp/out'
cp "${BASEDIR}/icon.icns" '/tmp/out/icon.icns'

# Build app
export GOOS="darwin"
export BINARIES=""
for ARCH in ${ARCHITECTURES}; do
  export GOARCH="${ARCH}"

  # See https://github.com/pojntfx/bagop/blob/main/main.go#L45
  export DEBARCH="${GOARCH}"
  export MACPORTS_ARGS=""
  if [ "${ARCH}" = "386" ]; then
    export DEBARCH="i686"
  elif [ "${ARCH}" = "amd64" ]; then
    export DEBARCH="x86_64"
  elif [ "${ARCH}" = "arm" ]; then
    export DEBARCH="armv7l"
  elif [ "${ARCH}" = "arm64" ]; then
    export MACPORTS_ARGS="--arm64"
    export DEBARCH="aarch64"
  fi

  export CC="$(find /osxcross/bin/${DEBARCH}-apple-darwin*-cc)"
  export CXX="$(find /osxcross/bin/${DEBARCH}-apple-darwin*-c++)"
  export PKGCONFIG="$(find /osxcross/bin/${DEBARCH}-apple-darwin*-pkg-config)"
  export PKG_CONFIG="$(find /osxcross/bin/${DEBARCH}-apple-darwin*-pkg-config)"

  # Install MacPorts packages
  if [ "${MACPORTS}" != "" ]; then
    rm -rf /osxcross/macports/pkgs/
    osxcross-macports install ${MACPORTS_ARGS} --static ${MACPORTS}
  fi

  go build -ldflags="-X github.com/pojntfx/hydrapp/hydrapp/pkg/ui.SelfUpdaterBranchTimestampRFC3339=${BRANCH_TIMESTAMP_RFC3339} -X github.com/pojntfx/hydrapp/hydrapp/pkg/ui.SelfUpdaterBranchID=${BRANCH_ID} -X github.com/pojntfx/hydrapp/hydrapp/pkg/ui.SelfUpdaterPackageType=dmg" -o "/tmp/${APP_ID}.${GOOS}-${DEBARCH}" "${GOMAIN}"

  export BINARIES="${BINARIES} /tmp/${APP_ID}.${GOOS}-${DEBARCH}"
done

lipo -create -output "/tmp/out/${APP_ID}.${GOOS}" ${BINARIES}

cp "${BASEDIR}/Info.plist" '/tmp/out/'

mkdir -p "/tmp/out/${APP_ID}.${GOOS}.dmg.mnt/${APP_NAME}.app/Contents/"{MacOS,Resources}
cp "/tmp/out/${APP_ID}.${GOOS}" "/tmp/out/${APP_ID}.${GOOS}.dmg.mnt/${APP_NAME}.app/Contents/MacOS/${APP_ID}"
cp '/tmp/out/Info.plist' "/tmp/out/${APP_ID}.${GOOS}.dmg.mnt/${APP_NAME}.app/Contents"
cp '/tmp/out/icon.icns' "/tmp/out/${APP_ID}.${GOOS}.dmg.mnt/${APP_NAME}.app/Contents/Resources"
rcodesign sign "/tmp/out/${APP_ID}.${GOOS}.dmg.mnt/${APP_NAME}.app"

genisoimage -V "Install $(echo ${APP_NAME} | cut -c -24)" -D -R -apple -no-pad -o "/tmp/out/${APP_ID}.${GOOS}-uncompressed.dmg" "/tmp/out/${APP_ID}.${GOOS}.dmg.mnt"
dmg "/tmp/out/${APP_ID}.${GOOS}-uncompressed.dmg" "/tmp/out/${APP_ID}.${GOOS}.dmg"
rcodesign sign "/tmp/out/${APP_ID}.${GOOS}.dmg"
gpg --detach-sign --armor "/tmp/out/${APP_ID}.${GOOS}.dmg"

cp "/tmp/out/${APP_ID}.${GOOS}.dmg" "/tmp/out/${APP_ID}.${GOOS}.dmg.asc" "/hydrapp/dst"

cd /hydrapp/dst

gpg --output "repo.asc" --armor --export
tree -J . -I 'index.html|index.json' | jq '.[0].contents' | jq ". |= map( . + {time: \"${BRANCH_TIMESTAMP_RFC3339}\"} )" | tee 'index.json'

if [ "${DST_UID}" != "" ] && [ "${DST_GID}" != "" ]; then
  chown -R "${DST_UID}:${DST_GID}" /hydrapp/dst
fi
