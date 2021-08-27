all: build

check:
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

build: check
	mkdir -p out
	bash -O extglob -c 'cd out && rm -rf -- !(*.jar)'
	CGO_ENABLED=1 GOOS=android GOARCH=arm64 CC="$$ANDROID_HOME/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang" CXX="$$ANDROID_HOME/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang++" go build -buildmode=c-shared -o=out/lib/arm64-v8a/libbackend.so ./pkg/main/;\
	cd out || exit 1;\
	cp ../*.java ../AndroidManifest.xml .;\
	javac -source 1.8 -target 1.8 -cp *.jar -bootclasspath $${ANDROID_HOME}/platforms/android-$${ANDROID_API_VERSION}/android.jar *.java;\
	$$ANDROID_HOME/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/d8 *.jar *.class --release;\
	$$ANDROID_HOME/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/aapt2 link -o $${APP_ID}.unsigned.apk -I *.jar -I $$ANDROID_HOME/platforms/android-$${ANDROID_API_VERSION}/android.jar --manifest AndroidManifest.xml;\
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
	mkdir -p out
	curl -Lo /tmp/Android-AdvancedWebView-v3.2.1.aar https://jitpack.io/com/github/delight-im/Android-AdvancedWebView/v3.2.1/Android-AdvancedWebView-v3.2.1.aar
	unzip -p /tmp/Android-AdvancedWebView-v3.2.1.aar classes.jar > out/Android-AdvancedWebView-v3.2.1.jar
