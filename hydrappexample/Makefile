DESTDIR ?=
PREFIX ?= /usr/local

SIZES ?= 16x16 22x22 24x24 32x32 36x36 48x48 64x64 72x72 96x96 128x128 192x192 256x256 512x512

all: build

build:
	go build -o out/com.pojtinger.felicitas.hydrappexample
	for icon in $(SIZES); do convert icon.png -resize $${icon} out/icon-$${icon}.png; done

install: build
	install -D -m 0755 out/com.pojtinger.felicitas.hydrappexample $(DESTDIR)$(PREFIX)/bin/com.pojtinger.felicitas.hydrappexample
	desktop-file-install --dir=$(DESTDIR)$(PREFIX)/share/applications com.pojtinger.felicitas.hydrappexample.desktop
	for icon in $(SIZES); do install -D -m 0644 out/icon-$${icon}.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/$${icon}/apps/com.pojtinger.felicitas.hydrappexample.png; done

uninstall:
	rm $(DESTDIR)$(PREFIX)/bin/com.pojtinger.felicitas.hydrappexample
	rm $(DESTDIR)$(PREFIX)/share/applications/com.pojtinger.felicitas.hydrappexample.desktop
	for icon in $(SIZES); do rm $(DESTDIR)$(PREFIX)/share/icons/hicolor/$${icon}/apps/com.pojtinger.felicitas.hydrappexample.png; done

run: build
	out/com.pojtinger.felicitas.hydrappexample

clean:
	rm -rf out .flatpak-builder build-dir repo

check-android:
# export APP_ID=com.pojtinger.felicitas.hydrappexample && export CERT_PASSWORD=123456 && export ANDROID_BUILD_TOOLS_VERSION=31.0.0 && export ANDROID_API_VERSION=30 && export ANDROID_NDK_VERSION=30 && export ANDROID_LEGACY_HOME="${ANDROID_HOME}/ndk/android-ndk-r10e" && export ANDROID_LEGACY_API_VERSION=8 && export ANDROID_LEGACY_NDK_VERSION=16
ifndef APP_ID
	$(error APP_ID env variable is undefined)
endif
ifndef CERT_PASSWORD
	$(error CERT_PASSWORD env variable is undefined)
endif
ifndef ANDROID_BUILD_TOOLS_VERSION
	$(error ANDROID_BUILD_TOOLS_VERSION env variable is undefined)
endif
ifndef ANDROID_API_VERSION
	$(error ANDROID_API_VERSION env variable is undefined)
endif
ifndef ANDROID_NDK_VERSION
	$(error ANDROID_NDK_VERSION env variable is undefined)
endif
ifndef ANDROID_HOME
	$(error ANDROID_HOME env variable is undefined)
endif
ifndef ANDROID_LEGACY_HOME
	$(error ANDROID_LEGACY_HOME env variable is undefined)
endif
ifndef ANDROID_LEGACY_API_VERSION
	$(error ANDROID_LEGACY_API_VERSION env variable is undefined)
endif
ifndef ANDROID_LEGACY_NDK_VERSION
	$(error ANDROID_LEGACY_NDK_VERSION env variable is undefined)
endif

build-android: check-android
	mkdir -p out
	bash -O extglob -c 'cd out && rm -rf -- !(*.jar)'
	mkdir -p out/drawable
	CGO_ENABLED=1 GOOS=android GOARCH=386 CC="$${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/i686-linux-android$${ANDROID_NDK_VERSION}-clang" CXX="$${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/i686-linux-android$${ANDROID_NDK_VERSION}-clang++" go build -buildmode=c-shared -o=out/lib/x86/libbackend.so
	CGO_ENABLED=1 GOOS=android GOARCH=amd64 CC="$${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android$${ANDROID_NDK_VERSION}-clang" CXX="$${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android$${ANDROID_NDK_VERSION}-clang++" go build -buildmode=c-shared -o=out/lib/x86_64/libbackend.so
	CGO_ENABLED=1 GOOS=android GOARCH=arm GOARM=5 CGO_LDFLAGS="--sysroot $${ANDROID_LEGACY_HOME}/platforms/android-$${ANDROID_LEGACY_API_VERSION}/arch-arm" CGO_CFLAGS="--sysroot $${ANDROID_LEGACY_HOME}/platforms/android-$${ANDROID_LEGACY_API_VERSION}/arch-arm" CC="$${ANDROID_LEGACY_HOME}/toolchains/arm-linux-androideabi-4.9/prebuilt/linux-x86_64/bin/arm-linux-androideabi-gcc" CXX="$${ANDROID_LEGACY_HOME}/toolchains/arm-linux-androideabi-4.9/prebuilt/linux-x86_64/bin/arm-linux-androideabi-g++" go build -tags "netgo,androiddnsfix,tlscertembed" -buildmode=c-shared -o=out/lib/armeabi/libbackend_compat.so
	CGO_ENABLED=1 GOOS=android GOARCH=arm GOARM=5 CC="$${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi$${ANDROID_LEGACY_NDK_VERSION}-clang" CXX="$${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi$${ANDROID_LEGACY_NDK_VERSION}-clang++" go build -tags "netgo,androiddnsfix,tlscertembed" -buildmode=c-shared -o=out/lib/armeabi-v7a/libbackend_compat.so
	CGO_ENABLED=1 GOOS=android GOARCH=arm CC="$${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi$${ANDROID_API_VERSION}-clang" CXX="$${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi$${ANDROID_API_VERSION}-clang++" go build -buildmode=c-shared -o=out/lib/armeabi-v7a/libbackend.so
	CGO_ENABLED=1 GOOS=android GOARCH=arm64 CC="$${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android$${ANDROID_NDK_VERSION}-clang" CXX="$${ANDROID_HOME}/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android$${ANDROID_NDK_VERSION}-clang++" go build -buildmode=c-shared -o=out/lib/arm64-v8a/libbackend.so
	cd out || exit 1;\
	cp ../*.java ../AndroidManifest.xml .;\
	cp ../icon.png drawable;\
	javac -source 1.8 -target 1.8 -cp *.jar -cp $${ANDROID_HOME}/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/core-lambda-stubs.jar -bootclasspath $${ANDROID_HOME}/platforms/android-$${ANDROID_API_VERSION}/android.jar *.java;\
	$${ANDROID_HOME}/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/d8 *.class --release; # Append *.jar here to use an external library\
	$${ANDROID_HOME}/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/aapt2 compile drawable/icon.png -o .;\
	$${ANDROID_HOME}/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/aapt2 link -o $${APP_ID}.unsigned.apk -I $${ANDROID_HOME}/platforms/android-$${ANDROID_API_VERSION}/android.jar *.flat --manifest AndroidManifest.xml; # Append -I *.jar here to use an external library\
	zip -ur $${APP_ID}.unsigned.apk lib classes.dex;\
	mkdir -p ../.android-certs;\
	[ ! -f ../.android-certs/$${APP_ID}.keystore ] && keytool -genkeypair -validity 365 -keystore ../.android-certs/$${APP_ID}.keystore -keyalg RSA -keysize 2048 -keypass $${CERT_PASSWORD} -storepass $${CERT_PASSWORD} -dname CN=$${APP_ID};\
	$${ANDROID_HOME}/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/zipalign -f -p 4 $${APP_ID}.unsigned.apk $${APP_ID}.apk;\
	$${ANDROID_HOME}/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/apksigner sign --ks ../.android-certs/$${APP_ID}.keystore --ks-pass pass:$${CERT_PASSWORD} --key-pass pass:$${CERT_PASSWORD} $${APP_ID}.apk

install-android: build-android
	$${ANDROID_HOME}/platform-tools/adb install -r out/$${APP_ID}.apk

uninstall-android: check-android
	$${ANDROID_HOME}/platform-tools/adb uninstall $${APP_ID}

run-android: install-android
	$${ANDROID_HOME}/platform-tools/adb shell am start -n $${APP_ID}/$${APP_ID}.MainActivity

clean-android:
	rm -rf out

depend-android: clean-android
	# You can add external dependencies here
	# mkdir -p out
	# curl -Lo /tmp/Android-AdvancedWebView-v3.2.1.aar https://jitpack.io/com/github/delight-im/Android-AdvancedWebView/v3.2.1/Android-AdvancedWebView-v3.2.1.aar
	# unzip -p /tmp/Android-AdvancedWebView-v3.2.1.aar classes.jar > out/Android-AdvancedWebView-v3.2.1.jar
