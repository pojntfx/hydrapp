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
export COMMIT_TIME_RFC3339="$(git log -1 --date=format:'%Y-%m-%dT%H:%M:%SZ' --format=%cd)"
export BRANCH_ID="stable"
if [ "$(git tag --points-at HEAD)" = "" ]; then
    export BRANCH_ID="$(git symbolic-ref --short HEAD)"
fi

CGO_ENABLED=0 bagop -j "$(nproc)" -b "${APP_ID}" -x "${GOEXCLUDE}" -d /dst -p "go build -o \$DST -ldflags='-X github.com/pojntfx/hydrapp/hydrapp/pkg/update.CommitTimeRFC3339=${COMMIT_TIME_RFC3339} -X github.com/pojntfx/hydrapp/hydrapp/pkg/update.BranchID=${BRANCH_ID}' ${GOMAIN}"

for FILE in /dst/*; do
    gpg --detach-sign --armor "${FILE}"
done

cd /dst

tree -T "${APP_NAME}" --du -h -D -H . -I 'index.html|index.json' -o 'index.html'
tree -J . -I 'index.html|index.json' | jq '.[0].contents' | jq ". |= map( . + {time: \"${COMMIT_TIME_RFC3339}\"} )" | tee 'index.json'
