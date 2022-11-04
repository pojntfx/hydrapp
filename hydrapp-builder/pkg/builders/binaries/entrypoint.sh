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

# Install host packages
if [ "${HOST_PACKAGES}" != "" ]; then
    apt update
    apt install -y ${HOST_PACKAGES}
fi

# Generate dependencies
GOFLAGS="${GOFLAGS}" sh -c "${GOGENERATE}"

# Build
GOFLAGS="-tags=selfupdate ${GOFLAGS}" CGO_ENABLED=0 bagop -j "$(nproc)" -b "${APP_ID}" -x "${GOEXCLUDE}" -d /dst "${GOMAIN}"

for FILE in /dst/*; do
    gpg --detach-sign --armor "${FILE}"
done

cd /dst

tree --timefmt '%Y-%m-%dT%H:%M:%SZ' -T "${APP_NAME}" --du -h -D -H . -I 'index.html|index.json' -o 'index.html'
tree --timefmt '%Y-%m-%dT%H:%M:%SZ' -J . -I 'index.html|index.json' | jq '.[0].contents' | tee 'index.json'
