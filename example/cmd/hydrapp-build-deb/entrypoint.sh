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

dpkg-source -b .

export BASEDIR="${PWD}"

export IFS='@'
for TARGET in ${TARGETS}; do
    cd "${BASEDIR}"

    export OS="$(cut -d'|' -f1 <<<"${TARGET}")"
    export DIST="$(cut -d'|' -f2 <<<"${TARGET}")"
    export MIRRORSITE="$(cut -d'|' -f3 <<<"${TARGET}")"
    export COMPONENTS="$(cut -d'|' -f4 <<<"${TARGET}")"
    export DEBOOTSTRAPOPTS="$(cut -d'|' -f5 <<<"${TARGET}")"
    export ARCHITECTURES="$(cut -d'|' -f6 <<<"${TARGET}")"

    export IFS=' '
    for ARCH in ${ARCHITECTURES}; do
        export ARCH

        pbuilder --create --mirror "${MIRRORSITE}" --components "${COMPONENTS}" $([ "${DEBOOTSTRAPOPTS}" != "" ] && echo --debootstrapopts "${DEBOOTSTRAPOPTS}")
        pbuilder build --mirror "${MIRRORSITE}" --components "${COMPONENTS}" $([ "${DEBOOTSTRAPOPTS}" != "" ] && echo --debootstrapopts "${DEBOOTSTRAPOPTS}") "../${PACKAGE}.dsc"

        for FILE in "/var/cache/pbuilder/${OS}-${DIST}-${ARCH}/result/"*; do
            mv "${FILE}" "${FILE//_${ARCH}/_${OS}_${DIST}_${ARCH}}" || :
        done

        mkdir -p "/dst/${OS}/pool/main"
        cp "/var/cache/pbuilder/${OS}-${DIST}-${ARCH}/result/"* "/dst/${OS}/pool/main" || :
    done

    cd '/dst' || exit 1

    for ARCH in ${ARCHITECTURES}; do
        mkdir -p "${OS}/dists/${DIST}/main/binary-${ARCH}"
    done

    mkdir -p "${OS}/dists/${DIST}/main/source" 'cache'

    cat >'apt-ftparchive.conf' <<EOT
Dir {
	ArchiveDir "./${OS}";
	CacheDir "./cache";
};
Default {
	Packages::Compress ". gzip bzip2";
	Sources::Compress ". gzip bzip2";
	Contents::Compress ". gzip bzip2";
};
TreeDefault {
	BinCacheDB "packages-\$(SECTION)-\$(ARCH).db";
	Directory "pool/\$(SECTION)";
	Packages "\$(DIST)/\$(SECTION)/binary-\$(ARCH)/Packages";
	SrcDirectory "pool/\$(SECTION)";
	Sources "\$(DIST)/\$(SECTION)/source/Sources";
	Contents "\$(DIST)/Contents-\$(ARCH)";
};
Tree "dists/${DIST}" {
	Sections "main";
	Architectures "${ARCHITECTURES} source";
}
EOT

    apt-ftparchive generate 'apt-ftparchive.conf'

    cat >"${OS}-${DIST}.conf" <<EOT
APT::FTPArchive::Release::Codename "${DIST}";
APT::FTPArchive::Release::Origin "Hydrapp APT repo";
APT::FTPArchive::Release::Components "main";
APT::FTPArchive::Release::Label "Packages for Hydrapp";
APT::FTPArchive::Release::Architectures "${ARCHITECTURES} source";
APT::FTPArchive::Release::Suite "${DIST}";
EOT

    apt-ftparchive -c "${OS}-${DIST}.conf" release "${OS}/dists/${DIST}" >"${OS}/dists/${DIST}/Release"

    gpg --output "repo.asc" --armor --export

    gpg --output "${OS}/dists/${DIST}/Release.gpg" -ba "${OS}/dists/${DIST}/Release"

    export IFS='@'
done
