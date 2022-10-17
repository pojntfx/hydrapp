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

echo "%_signature gpg
%_gpg_name ${GPG_KEY_ID}" >"${HOME}/.rpmmacros"

# Build tarball and source package
export PACKAGE="${APP_ID}-${PACKAGE_VERSION}"
export SUFFIX="${PACKAGE_SUFFIX}"
export SPEC="${APP_ID}.spec"

# See https://github.com/pojntfx/bagop/blob/main/main.go#L45
export DEBARCH="${ARCHITECTURE}"
if [ "${ARCHITECTURE}" = "386" ]; then
    export DEBARCH="i686"
elif [ "${ARCHITECTURE}" = "amd64" ]; then
    export DEBARCH="x86_64"
elif [ "${ARCHITECTURE}" = "arm" ]; then
    export DEBARCH="armv7l"
elif [ "${ARCHITECTURE}" = "arm64" ]; then
    export DEBARCH="aarch64"
fi

rpmdev-setuptree

export TARBALL="${HOME}/rpmbuild/SOURCES/${PACKAGE}.tar.gz"
export DSC="${HOME}/rpmbuild/SRPMS/${PACKAGE}-${SUFFIX}.src.rpm"

tar -cvzf "${TARBALL}" --exclude out --transform "s,^,${PACKAGE}/," .
rpmbuild -bs "${SPEC}"

rpmlint "${DSC}"

# Build chroot
mock -r "${DISTRO}-${DEBARCH}" "${DSC}" --enable-network

rpmlint "/var/lib/mock/${DISTRO}-${DEBARCH}/result/*.rpm"

mkdir -p '/dst'
cp "/var/lib/mock/${DISTRO}-${DEBARCH}/result/"*.rpm '/dst'

rpm --addsign '/dst'/*.rpm

createrepo '/dst'

gpg --detach-sign --armor "/dst/repodata/repomd.xml"

gpg --output "/dst/repodata/repo.asc" --armor --export

# Add repo file
echo "[hydrapp-repo]
name=Hydrapp YUM repo
baseurl=${BASE_URL}
enabled=1
gpgcheck=1
gpgkey=${BASE_URL}/repodata/repo.asc" >"/dst/repodata/hydrapp.repo"
