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
export COMMIT_TIME_UNIX="$(git log -1 --format=%ct)"

echo "%_signature gpg
%_gpg_name $(echo ${PGP_KEY_ID} | base64 -d)
%commit_time_unix ${COMMIT_TIME_UNIX}" >"${HOME}/.rpmmacros"

# Build tarball and source package
export PACKAGE="${APP_ID}-${PACKAGE_VERSION}"
export SUFFIX="${COMMIT_TIME_UNIX}.${PACKAGE_SUFFIX}"
export SPEC="${BASEDIR}/${APP_ID}.spec"

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
mock --define "commit_time_unix ${COMMIT_TIME_UNIX}" -r "${DISTRO}-${DEBARCH}" "${DSC}" --enable-network

rpmlint "/var/lib/mock/${DISTRO}-${DEBARCH}/result/*.rpm"

mkdir -p '/hydrapp/dst'
cp "/var/lib/mock/${DISTRO}-${DEBARCH}/result/"*.rpm '/hydrapp/dst'

rpm --addsign '/hydrapp/dst'/*.rpm

createrepo '/hydrapp/dst'

gpg --detach-sign --armor "/hydrapp/dst/repodata/repomd.xml"

gpg --output "/hydrapp/dst/repodata/repo.asc" --armor --export

# Add repo file
echo "[hydrapp-repo]
name=hydrapp YUM repo
baseurl=${BASE_URL}
enabled=1
gpgcheck=1
gpgkey=${BASE_URL}/repodata/repo.asc" >"/hydrapp/dst/repodata/hydrapp.repo"

if [ "${DST_UID}" != "" ] && [ "${DST_GID}" != "" ]; then
    chown -R "${DST_UID}:${DST_GID}" /hydrapp/dst
fi
