all: build

build: clean
	CGO_ENABLED=1 GOOS=android GOARCH=arm64 CC="$$ANDROID_HOME/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang" CXX="$$ANDROID_HOME/ndk-bundle/toolchains/llvm/prebuilt/linux-x86_64/bin/aarch64-linux-android21-clang++" CGO_CFLAGS="-I "$$ANDROID_HOME/ndk-bundle/sources/android/native_app_glue"" go build -buildmode=c-shared -o=out/lib/arm64-v8a/libnativeactivity.so ./pkg/main/;\
	cd out || exit 1;\
	cp ../*.java ../AndroidManifest.xml .;\
	javac -source 1.8 -target 1.8 -cp *.jar -bootclasspath $${ANDROID_HOME}/platforms/android-$${ANDROID_API_VERSION}/android.jar *.java;\
	$$ANDROID_HOME/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/d8 *.class --release;\
	$$ANDROID_HOME/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/aapt2 link -o $${APP_ID}.unsigned.apk -I $$ANDROID_HOME/platforms/android-$${ANDROID_API_VERSION}/android.jar --manifest AndroidManifest.xml;\
	zip -ur $${APP_ID}.unsigned.apk lib classes.dex;\
	mkdir -p ~/.android-certs;\
	[ ! -f ~/.android-certs/$${APP_ID}.keystore ] && keytool -genkeypair -validity 365 -keystore ~/.android-certs/$${APP_ID}.keystore -keyalg RSA -keysize 2048 -keypass $${CERT_PASSWORD} -storepass $${CERT_PASSWORD} -dname CN=$${APP_ID};\
	$$ANDROID_HOME/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/zipalign -f -p 4 $${APP_ID}.unsigned.apk $${APP_ID}.apk;\
	$$ANDROID_HOME/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/apksigner sign --ks ~/.android-certs/$${APP_ID}.keystore --ks-pass pass:$${CERT_PASSWORD} --key-pass pass:$${CERT_PASSWORD} $${APP_ID}.apk;\

install: build
	$$ANDROID_HOME/platform-tools/adb install out/$${APP_ID}.apk

run: install
	$$ANDROID_HOME/platform-tools/adb shell am start -n $${APP_ID}/com.example.app.MainActivity

clean:
	rm -rf out