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

# Build chroot and source package
export PACKAGE="${APP_ID}_${PACKAGE_VERSION}"
