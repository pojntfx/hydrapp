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

mkdir -p '/tmp/out'
bash -O extglob -c 'cd /tmp/out && rm -rf -- !(*.jar)'
mkdir -p '/tmp/out/drawable'

# Build native libraries
CGO_ENABLED=1 GOOS=android GOARCH=386 CC="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-$(uname -m)/bin/i686-linux-android${ANDROID_NDK_VERSION}-clang" CXX="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-$(uname -m)/bin/i686-linux-android${ANDROID_NDK_VERSION}-clang++" go build -buildmode='c-shared' -o='/tmp/out/lib/x86/libbackend.so' "${GOMAIN}"
CGO_ENABLED=1 GOOS=android GOARCH=amd64 CC="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-$(uname -m)/bin/$(uname -m)-linux-android${ANDROID_NDK_VERSION}-clang" CXX="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-$(uname -m)/bin/$(uname -m)-linux-android${ANDROID_NDK_VERSION}-clang++" go build -buildmode='c-shared' -o='/tmp/out/lib/$(uname -m)/libbackend.so' "${GOMAIN}"
# Only build with the legacy Android NDK if we're on x86_64
if [ "$(uname -m)" = "x86_64" ]; then
  CGO_ENABLED=1 GOOS=android GOARCH=arm GOARM=5 CGO_LDFLAGS="--sysroot ${ANDROID_LEGACY_HOME}/platforms/android-${ANDROID_LEGACY_API_VERSION}/arch-arm" CGO_CFLAGS="--sysroot ${ANDROID_LEGACY_HOME}/platforms/android-${ANDROID_LEGACY_API_VERSION}/arch-arm" CC="${ANDROID_LEGACY_HOME}/toolchains/arm-linux-androideabi-4.9/prebuilt/linux-$(uname -m)/bin/arm-linux-androideabi-gcc" CXX="${ANDROID_LEGACY_HOME}/toolchains/arm-linux-androideabi-4.9/prebuilt/linux-$(uname -m)/bin/arm-linux-androideabi-g++" go build -tags "netgo,androiddnsfix,tlscertembed" -buildmode='c-shared' -o='/tmp/out/lib/armeabi/libbackend_compat.so' "${GOMAIN}"
  CGO_ENABLED=1 GOOS=android GOARCH=arm GOARM=5 CC="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-$(uname -m)/bin/armv7a-linux-androideabi${ANDROID_LEGACY_NDK_VERSION}-clang" CXX="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-$(uname -m)/bin/armv7a-linux-androideabi${ANDROID_LEGACY_NDK_VERSION}-clang++" go build -tags "netgo,androiddnsfix,tlscertembed" -buildmode='c-shared' -o='/tmp/out/lib/armeabi-v7a/libbackend_compat.so' "${GOMAIN}"
fi
CGO_ENABLED=1 GOOS=android GOARCH=arm CC="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-$(uname -m)/bin/armv7a-linux-androideabi${ANDROID_NDK_VERSION}-clang" CXX="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-$(uname -m)/bin/armv7a-linux-androideabi${ANDROID_API_VERSION}-clang++" go build -buildmode='c-shared' -o='/tmp/out/lib/armeabi-v7a/libbackend.so' "${GOMAIN}"
CGO_ENABLED=1 GOOS=android GOARCH=arm64 CC="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-$(uname -m)/bin/aarch64-linux-android${ANDROID_NDK_VERSION}-clang" CXX="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-$(uname -m)/bin/aarch64-linux-android${ANDROID_NDK_VERSION}-clang++" go build -buildmode='c-shared' -o='/tmp/out/lib/arm64-v8a/libbackend.so' "${GOMAIN}"

# Sign native libraries with PGP
gpg --detach-sign --armor "/tmp/out/lib/"*/*

# Create package
cd '/tmp/out' || exit 1
cp "${BASEDIR}"/*.java "${BASEDIR}/AndroidManifest.xml" .
cp "${BASEDIR}/icon.png" 'drawable'
javac -source "1.8" -target "1.8" -cp *.jar -cp "${ANDROID_HOME}/build-tools/${ANDROID_BUILD_TOOLS_VERSION}/core-lambda-stubs.jar" -bootclasspath "${ANDROID_HOME}/platforms/android-${ANDROID_API_VERSION}/android.jar" *.java
"${ANDROID_HOME}/build-tools/${ANDROID_BUILD_TOOLS_VERSION}/d8" *.class --release
"${ANDROID_HOME}/build-tools/${ANDROID_BUILD_TOOLS_VERSION}/aapt2" compile 'drawable/icon.png' -o .
"${ANDROID_HOME}/build-tools/${ANDROID_BUILD_TOOLS_VERSION}/aapt2" link -o "${APP_ID}.unsigned" -I "${ANDROID_HOME}/platforms/android-${ANDROID_API_VERSION}/android.jar" *.flat --manifest 'AndroidManifest.xml'
zip -ur "${APP_ID}.unsigned" 'lib' 'classes.dex'
mkdir -p "/tmp/out/android-certs" # Append *.jar here to use an external library

# Sign package with Android certificate
echo "${JAVA_KEYSTORE}" | base64 -d >"/tmp/out/android-certs/${APP_ID}.keystore"

export ANDROID_CERT_CN="$(keytool -noprompt -storepass $(echo ${JAVA_KEYSTORE_PASSWORD} | base64 -d) -keypass $(echo ${JAVA_CERTIFICATE_PASSWORD} | base64 -d) -v -list -keystore /tmp/out/android-certs/${APP_ID}.keystore | grep -oP 'Owner: CN=\K\w(.*)')"
export ANDROID_CERT_ALIAS="$(keytool -noprompt -storepass $(echo ${JAVA_KEYSTORE_PASSWORD} | base64 -d) -keypass $(echo ${JAVA_CERTIFICATE_PASSWORD} | base64 -d) -v -list -keystore /tmp/out/android-certs/${APP_ID}.keystore | grep -oP 'Alias name: \K\w(.*)')"

"${ANDROID_HOME}/build-tools/${ANDROID_BUILD_TOOLS_VERSION}/zipalign" -f -p 4 "${APP_ID}.unsigned" "${APP_ID}.apk"
"${ANDROID_HOME}/build-tools/${ANDROID_BUILD_TOOLS_VERSION}/apksigner" sign --ks "/tmp/out/android-certs/${APP_ID}.keystore" --ks-pass pass:"$(echo ${JAVA_KEYSTORE_PASSWORD} | base64 -d)" --key-pass pass:"$(echo ${JAVA_CERTIFICATE_PASSWORD} | base64 -d)" "${APP_ID}.apk"

# Sign package with PGP and stage
gpg --detach-sign --armor "${APP_ID}.apk"

# Setup repository
rm -rf "/hydrapp/dst/"*
cd "/hydrapp/dst" || exit 1

fdroid init
cp -f ${BASEDIR}/config.yml config.yml
perl -p -i -e 's/\{ ANDROID_HOME \}/$ENV{"ANDROID_HOME"}/g' config.yml
perl -p -i -e 's/\{ ANDROID_CERT_ALIAS \}/$ENV{"ANDROID_CERT_ALIAS"}/g' config.yml
perl -p -i -e 's/\{ JAVA_KEYSTORE_PASSWORD \}/`echo $ENV{"JAVA_KEYSTORE_PASSWORD"} | base64 -d`/ge' config.yml
perl -p -i -e 's/\{ JAVA_CERTIFICATE_PASSWORD \}/`echo $ENV{"JAVA_CERTIFICATE_PASSWORD"} | base64 -d`/ge' config.yml
perl -p -i -e 's/\{ ANDROID_CERT_CN \}/$ENV{"ANDROID_CERT_CN"}/g' config.yml

cp "/tmp/out/${APP_ID}.apk" 'repo/'
cp "${BASEDIR}/icon.png" .
cp "/tmp/out/android-certs/${APP_ID}.keystore" 'keystore.p12'

fdroid update --create-metadata
fdroid gpgsign

cat >'.gitignore' <<'EOT'
*.ks
*.jks
*.keystore
*.crt
*config.py
tmp
EOT

if [ "${DST_UID}" != "" ] && [ "${DST_GID}" != "" ]; then
  chown -R "${DST_UID}:${DST_GID}" /hydrapp/dst
fi
