all: build

check:
# export APP_ID=com.pojtinger.felicitas.integratedWebserverExample && export CERT_PASSWORD=123456 && export ANDROID_BUILD_TOOLS_VERSION=31.0.0 && export ANDROID_API_VERSION=30 && export ANDROID_LEGACY_HOME=$ANDROID_HOME/ndk/android-ndk-r10e && export ANDROID_LEGACY_API_VERSION=8
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
ifndef ANDROID_HOME
	$(error ANDROID_HOME env variable is undefined)
endif
ifndef ANDROID_LEGACY_HOME
	$(error ANDROID_LEGACY_HOME env variable is undefined)
endif
ifndef ANDROID_LEGACY_API_VERSION
	$(error ANDROID_LEGACY_API_VERSION env variable is undefined)
endif

build: check
	mkdir -p out
	bash -O extglob -c 'cd out && rm -rf -- !(*.jar)'
	CGO_ENABLED=1 GOOS=android GOARCH=386 CC="$$ANDROID_HOME/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/i686-linux-android$${ANDROID_API_VERSION}-clang" CXX="$$ANDROID_HOME/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/i686-linux-android$${ANDROID_API_VERSION}-clang++" go build -buildmode=c-shared -o=out/lib/x86/libbackend.so ./pkg/main/;\
	CGO_ENABLED=1 GOOS=android GOARCH=amd64 CC="$$ANDROID_HOME/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android$${ANDROID_API_VERSION}-clang" CXX="$$ANDROID_HOME/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/x86_64-linux-android$${ANDROID_API_VERSION}-clang++" go build -buildmode=c-shared -o=out/lib/x86_64/libbackend.so ./pkg/main/;\
	CGO_ENABLED=1 GOOS=android GOARCH=arm GOARM=5 CGO_LDFLAGS="--sysroot $$ANDROID_LEGACY_HOME/platforms/android-$$ANDROID_LEGACY_API_VERSION/arch-arm" CGO_CFLAGS="--sysroot $$ANDROID_LEGACY_HOME/platforms/android-$$ANDROID_LEGACY_API_VERSION/arch-arm" CC="$$ANDROID_LEGACY_HOME/toolchains/arm-linux-androideabi-4.9/prebuilt/linux-x86_64/bin/arm-linux-androideabi-gcc" CXX="$$ANDROID_LEGACY_HOME/toolchains/arm-linux-androideabi-4.9/prebuilt/linux-x86_64/bin/arm-linux-androideabi-g++" go build -tags "netgo,androiddnsfix,tlscertembed" -buildmode=c-shared -o=out/lib/armeabi/libbackend.so ./pkg/main/;\
	CGO_ENABLED=1 GOOS=android GOARCH=arm CC="$$ANDROID_HOME/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi$${ANDROID_API_VERSION}-clang" CXX="$$ANDROID_HOME/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/armv7a-linux-androideabi$${ANDROID_API_VERSION}-clang++" go build -buildmode=c-shared -o=out/lib/armeabi-v7a/libbackend.so ./pkg/main/;\
	CGO_ENABLED=1 GOOS=android GOARCH=arm64 CC="$$ANDROID_HOME/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android$${ANDROID_API_VERSION}-clang" CXX="$$ANDROID_HOME/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android$${ANDROID_API_VERSION}-clang++" go build -buildmode=c-shared -o=out/lib/arm64-v8a/libbackend.so ./pkg/main/;\
	cd out || exit 1;\
	cp ../*.java ../AndroidManifest.xml .;\
	javac -source 1.8 -target 1.8 -cp *.jar -cp $$ANDROID_HOME/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/core-lambda-stubs.jar -bootclasspath $${ANDROID_HOME}/platforms/android-$${ANDROID_API_VERSION}/android.jar *.java;\
	$$ANDROID_HOME/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/d8 *.class --release; # Append *.jar here to use an external library\
	$$ANDROID_HOME/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/aapt2 link -o $${APP_ID}.unsigned.apk -I $$ANDROID_HOME/platforms/android-$${ANDROID_API_VERSION}/android.jar --manifest AndroidManifest.xml; # Append -I *.jar here to use an external library\
	zip -ur $${APP_ID}.unsigned.apk lib classes.dex;\
	mkdir -p ~/.android-certs;\
	[ ! -f ~/.android-certs/$${APP_ID}.keystore ] && keytool -genkeypair -validity 365 -keystore ~/.android-certs/$${APP_ID}.keystore -keyalg RSA -keysize 2048 -keypass $${CERT_PASSWORD} -storepass $${CERT_PASSWORD} -dname CN=$${APP_ID};\
	$$ANDROID_HOME/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/zipalign -f -p 4 $${APP_ID}.unsigned.apk $${APP_ID}.apk;\
	$$ANDROID_HOME/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/apksigner sign --ks ~/.android-certs/$${APP_ID}.keystore --ks-pass pass:$${CERT_PASSWORD} --key-pass pass:$${CERT_PASSWORD} $${APP_ID}.apk;\

install: build
	$$ANDROID_HOME/platform-tools/adb install out/$${APP_ID}.apk

run: install
	$$ANDROID_HOME/platform-tools/adb shell am start -n $${APP_ID}/$${APP_ID}.MainActivity

clean:
	rm -rf out

depend: clean
	# You can add external dependencies here
	# mkdir -p out
	# curl -Lo /tmp/Android-AdvancedWebView-v3.2.1.aar https://jitpack.io/com/github/delight-im/Android-AdvancedWebView/v3.2.1/Android-AdvancedWebView-v3.2.1.aar
	# unzip -p /tmp/Android-AdvancedWebView-v3.2.1.aar classes.jar > out/Android-AdvancedWebView-v3.2.1.jar
