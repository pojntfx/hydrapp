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

# Build chroot and source package
export PACKAGE="${APP_ID}_${PACKAGE_VERSION}"

if [ ! -d debian ]; then
	cp -r "${BASEDIR}/debian" debian
fi

dpkg-source -b .

pbuilder --create --mirror "${MIRRORSITE}" --components "${COMPONENTS}" $([ "${DEBOOTSTRAPOPTS}" != "" ] && echo --debootstrapopts "${DEBOOTSTRAPOPTS}")
pbuilder build --mirror "${MIRRORSITE}" --components "${COMPONENTS}" $([ "${DEBOOTSTRAPOPTS}" != "" ] && echo --debootstrapopts "${DEBOOTSTRAPOPTS}") "../${PACKAGE}.dsc"

for FILE in "/var/cache/pbuilder/${OS}-${DISTRO}-${ARCHITECTURE}/result/"*; do
	mv "${FILE}" "${FILE//_${ARCHITECTURE}/_${OS}_${DISTRO}_${ARCHITECTURE}}" || :
done

mkdir -p "/dst/pool/main"
cp "/var/cache/pbuilder/${OS}-${DISTRO}-${ARCHITECTURE}/result/"* "/dst/pool/main" || :

cd '/dst' || exit 1

mkdir -p "main/binary-${ARCHITECTURE}"

mkdir -p "main/source" 'cache'

cat >'apt-ftparchive.conf' <<EOT
Dir {
	ArchiveDir "${OS}";
	CacheDir "cache";
};
Default {
	Packages::Compress ". gzip bzip2";
	Sources::Compress ". gzip bzip2";
	Contents::Compress ". gzip bzip2";
};
TreeDefault {
	BinCacheDB "packages-\$(SECTION)-\$(ARCHITECTURE).db";
	Directory "pool/\$(SECTION)";
	Packages "\$(DISTRO)/\$(SECTION)/binary-\$(ARCHITECTURE)/Packages";
	SrcDirectory "pool/\$(SECTION)";
	Sources "\$(DISTRO)/\$(SECTION)/source/Sources";
	Contents "\$(DISTRO)/Contents-\$(ARCHITECTURE)";
};
Tree "." {
	Sections "main";
	ARCHITECTURE "${ARCHITECTURE} source";
}
EOT

apt-ftparchive generate 'apt-ftparchive.conf'

cat >"${OS}-${DISTRO}.conf" <<EOT
APT::FTPArchive::Release::Codename "${DISTRO}";
APT::FTPArchive::Release::Origin "hydrapp APT repo";
APT::FTPArchive::Release::Components "main";
APT::FTPArchive::Release::Label "Packages for hydrapp";
APT::FTPArchive::Release::ARCHITECTURE "${ARCHITECTURE} source";
APT::FTPArchive::Release::Suite "${DISTRO}";
EOT

apt-ftparchive -c "${OS}-${DISTRO}.conf" release "." >"Release"

gpg --output "repo.asc" --armor --export

gpg --output "Release.gpg" -ba "Release"

if [ "${DST_UID}" != "" ] && [ "${DST_GID}" != "" ]; then
	chown -R "${DST_UID}:${DST_GID}" /dst
fi
