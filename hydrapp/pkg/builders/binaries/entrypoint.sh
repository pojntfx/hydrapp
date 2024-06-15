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

# Install host packages
if [ "${HOST_PACKAGES}" != "" ]; then
    apt update
    apt install -y ${HOST_PACKAGES}
fi

# Generate dependencies
GOFLAGS="${GOFLAGS}" sh -c "${GOGENERATE}"

# Build
CGO_ENABLED=0 bagop -j "$(nproc)" -b "${APP_ID}" -x "${GOEXCLUDE}" -d /hydrapp/dst -p "go build -o \$DST -ldflags='-X github.com/pojntfx/hydrapp/hydrapp/pkg/ui.SelfUpdaterBranchTimestampRFC3339=${BRANCH_TIMESTAMP_RFC3339} -X github.com/pojntfx/hydrapp/hydrapp/pkg/ui.SelfUpdaterBranchID=${BRANCH_ID}' ${GOMAIN}"

for FILE in /hydrapp/dst/*; do
    gpg --detach-sign --armor "${FILE}"
done

cd /hydrapp/dst

gpg --output "repo.asc" --armor --export
tree -T "${APP_NAME}" --du -h -D -H . -I 'index.html|index.json' -o 'index.html'
tree -J . -I 'index.html|index.json' | jq '.[0].contents' | jq ". |= map( . + {time: \"${BRANCH_TIMESTAMP_RFC3339}\"} )" | tee 'index.json'

if [ "${DST_UID}" != "" ] && [ "${DST_GID}" != "" ]; then
    chown -R "${DST_UID}:${DST_GID}" /hydrapp/dst
fi
