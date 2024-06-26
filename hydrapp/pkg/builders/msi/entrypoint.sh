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

# Re-initialize WINE at runtime
rm -rf /root/.wine
wine64 cmd.exe /c dir
mv /tmp/msys64 /root/.wine/drive_c/msys64
mkdir -p '/root/.wine/drive_c/go' '/root/.wine/drive_c/tmp' '/root/.wine/drive_c/users/root/Documents/go-workspace'

# Prepare build environment
export BASEDIR="${PWD}/${GOMAIN}"

# Configure Go
export GOPROXY='https://proxy.golang.org,direct'

# Install MSYS2 packages
if [ "${MSYS2PACKAGES}" != "" ]; then
  wine64 bash.exe -c "pacman --noconfirm --ignore pacman --needed -S ${MSYS2PACKAGES}"
fi

# Generate dependencies
GOFLAGS="${GOFLAGS}" sh -c "${GOGENERATE}"

# Create icons
mkdir -p '/tmp/out'
cp "${BASEDIR}/icon.ico" '/tmp/out/icon.ico'

# Build app
export GOOS="windows"
export GOARCH="${ARCHITECTURE}"

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

if [ "${ARCHITECTURE}" = "amd64" ]; then
  GOPATH='/root/.wine/drive_c/go' go mod download -x

  cp -r . '/root/.wine/drive_c/users/root/Documents/go-workspace'
  rm -rf '/root/.wine/drive_c/users/root/Documents/go-workspace/out'

  wine64 bash.exe -c "export PATH=$PATH:/ucrt64/bin:/msys64/usr/bin GOPATH=/c/go GOROOT=/ucrt64/lib/go TMP=/c/tmp TEMP=/c/tmp GOARCH=amd64 CGO_ENABLED=1 GOPROXY='https://proxy.golang.org,direct' GOFLAGS=${GOFLAGS} && cd /c/users/root/Documents/go-workspace && git config --add safe.directory '*' && go build -ldflags='-linkmode=external -H=windowsgui -X github.com/pojntfx/hydrapp/hydrapp/pkg/ui.SelfUpdaterBranchTimestampRFC3339=${BRANCH_TIMESTAMP_RFC3339} -X github.com/pojntfx/hydrapp/hydrapp/pkg/ui.SelfUpdaterBranchID=${BRANCH_ID} -X github.com/pojntfx/hydrapp/hydrapp/pkg/ui.SelfUpdaterPackageType=msi' -x -v -o out/${APP_ID}.${GOOS}-${DEBARCH}.exe ${GOMAIN}"

  # Copy binaries to staging directory
  yes | cp -rf /root/.wine/drive_c/users/root/Documents/go-workspace/out/* '/tmp/out'

  find /root/.wine/drive_c/msys64/ucrt64/ -regex "${MSYS2INCLUDE}" -print0 | tar -c --null --files-from - | tar -C '/tmp/out' -x --strip-components=5
else
  go build -ldflags="-X github.com/pojntfx/hydrapp/hydrapp/pkg/ui.SelfUpdaterBranchTimestampRFC3339=${BRANCH_TIMESTAMP_RFC3339} -X github.com/pojntfx/hydrapp/hydrapp/pkg/ui.SelfUpdaterBranchID=${BRANCH_ID} -X github.com/pojntfx/hydrapp/hydrapp/pkg/ui.SelfUpdaterPackageType=msi" -o "/tmp/out/${APP_ID}.${GOOS}-${DEBARCH}.exe" "${GOMAIN}"
fi

cd '/tmp/out'

# Create and analyze files to include in the installer
find . -type f | wixl-heat -p ./ --directory-ref INSTALLDIR --component-group ApplicationContent --var 'var.SourceDir' >/tmp/hydrapp.wxi

xmllint --xpath "//*[local-name()='DirectoryRef']" /tmp/hydrapp.wxi >/tmp/hydrapp-directories.xml
xmllint --xpath "//*[local-name()='ComponentRef']" /tmp/hydrapp.wxi >/tmp/hydrapp-component-refs.xml

export STARTID="$(cat /tmp/hydrapp.wxi | grep ${APP_ID}.${GOOS}-${DEBARCH}.exe | xmllint --xpath 'string(//File/@Id)' -)"

# Build WiX installer
wixl -v -D SourceDir="." -v -o "/hydrapp/dst/${APP_ID}.${GOOS}-${DEBARCH}.msi" <(cat "${BASEDIR}/${APP_ID}.wxl" | perl -p -e 'use File::Slurp; my $text = read_file("/tmp/hydrapp-directories.xml"); s+<hydrappDirectories />+$text+g' | perl -p -e 'use File::Slurp; my $text = read_file("/tmp/hydrapp-component-refs.xml"); s+<hydrappComponentRefs />+$text+g' | perl -p -e 's+{ StartID }+$ENV{STARTID}+g')

gpg --detach-sign --armor "/hydrapp/dst/${APP_ID}.${GOOS}-${DEBARCH}.msi"

cd /hydrapp/dst

gpg --output "repo.asc" --armor --export
tree -J . -I 'index.html|index.json' | jq '.[0].contents' | jq ". |= map( . + {time: \"${BRANCH_TIMESTAMP_RFC3339}\"} )" | tee 'index.json'

if [ "${DST_UID}" != "" ] && [ "${DST_GID}" != "" ]; then
  chown -R "${DST_UID}:${DST_GID}" /hydrapp/dst
fi
