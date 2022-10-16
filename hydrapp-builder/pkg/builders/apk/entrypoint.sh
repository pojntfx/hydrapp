#!/bin/bash

set -e

# Setup workdir
mkdir -p /work
cp -r . /work
cd /work

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

# Prepare build environment
export BASEDIR="${PWD}"

# Generate dependencies
go generate ./...

mkdir -p '/tmp/out'
bash -O extglob -c 'cd /tmp/out && rm -rf -- !(*.jar)'
mkdir -p '/tmp/out/drawable'

# Build native libraries
CGO_ENABLED=1 GOOS=android GOARCH=386 CC="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/i686-linux-android${ANDROID_NDK_VERSION}-clang" CXX="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/i686-linux-android${ANDROID_NDK_VERSION}-clang++" go build -buildmode='c-shared' -o='/tmp/out/lib/x86/libbackend.so'
CGO_ENABLED=1 GOOS=android GOARCH=amd64 CC="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android${ANDROID_NDK_VERSION}-clang" CXX="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android${ANDROID_NDK_VERSION}-clang++" go build -buildmode='c-shared' -o='/tmp/out/lib/x86_64/libbackend.so'
CGO_ENABLED=1 GOOS=android GOARCH=arm GOARM=5 CGO_LDFLAGS="--sysroot ${ANDROID_LEGACY_HOME}/platforms/android-${ANDROID_LEGACY_API_VERSION}/arch-arm" CGO_CFLAGS="--sysroot ${ANDROID_LEGACY_HOME}/platforms/android-${ANDROID_LEGACY_API_VERSION}/arch-arm" CC="${ANDROID_LEGACY_HOME}/toolchains/arm-linux-androideabi-4.9/prebuilt/linux-x86_64/bin/arm-linux-androideabi-gcc" CXX="${ANDROID_LEGACY_HOME}/toolchains/arm-linux-androideabi-4.9/prebuilt/linux-x86_64/bin/arm-linux-androideabi-g++" go build -tags "netgo,androiddnsfix,tlscertembed" -buildmode='c-shared' -o='/tmp/out/lib/armeabi/libbackend_compat.so'
CGO_ENABLED=1 GOOS=android GOARCH=arm GOARM=5 CC="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi${ANDROID_LEGACY_NDK_VERSION}-clang" CXX="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi${ANDROID_LEGACY_NDK_VERSION}-clang++" go build -tags "netgo,androiddnsfix,tlscertembed" -buildmode='c-shared' -o='/tmp/out/lib/armeabi-v7a/libbackend_compat.so'
CGO_ENABLED=1 GOOS=android GOARCH=arm CC="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi${ANDROID_NDK_VERSION}-clang" CXX="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi${ANDROID_API_VERSION}-clang++" go build -buildmode='c-shared' -o='/tmp/out/lib/armeabi-v7a/libbackend.so'
CGO_ENABLED=1 GOOS=android GOARCH=arm64 CC="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android${ANDROID_NDK_VERSION}-clang" CXX="${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android${ANDROID_NDK_VERSION}-clang++" go build -buildmode='c-shared' -o='/tmp/out/lib/arm64-v8a/libbackend.so'

# Sign native libraries with GPG
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
echo "${ANDROID_CERT_CONTENT}" | base64 -d >"/tmp/out/android-certs/${APP_ID}.keystore"
"${ANDROID_HOME}/build-tools/${ANDROID_BUILD_TOOLS_VERSION}/zipalign" -f -p 4 "${APP_ID}.unsigned" "${APP_ID}.apk"
"${ANDROID_HOME}/build-tools/${ANDROID_BUILD_TOOLS_VERSION}/apksigner" sign --ks "/tmp/out/android-certs/${APP_ID}.keystore" --ks-pass pass:"$(echo ${ANDROID_CERT_PASSWORD} | base64 -d)" --key-pass pass:"$(echo ${ANDROID_CERT_PASSWORD} | base64 -d)" "${APP_ID}.apk"

# Sign package with GPG and stage
gpg --detach-sign --armor "${APP_ID}.apk"

# Setup repository
rm -rf "/dst/"*
cd "/dst" || exit 1

fdroid init
cat >'config.yml' <<EOT
---
sdk_path: "${ANDROID_HOME}"
repo_url: "${BASE_URL}"
repo_name: Hydrapp F-Droid Repo
repo_description: >-
  Android apps for Hydrapp.
repo_icon: icon.png
repo_keyalias: ${APP_ID}
keystore: keystore.p12
keystorepass: "$(echo ${ANDROID_CERT_PASSWORD} | base64 -d)"
keypass: "$(echo ${ANDROID_CERT_PASSWORD} | base64 -d)"
keydname: CN=${APP_ID}
apksigner: /usr/bin/apksigner
EOT

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
