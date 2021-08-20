all: build

build: clean
	mkdir -p build;\
	cd build;\
	cp ../*.java .;\
	gomobile bind -target android -androidapi 21 ../...;\
	javac -source 1.8 -target 1.8 -cp *.jar -bootclasspath $${ANDROID_HOME}/platforms/android-$${ANDROID_API_VERSION}/android.jar *.java;\
	for aar in *.aar; do unzip -p -j $${aar} classes.jar > $${aar}.classes.jar; unzip bindings.aar jni/*; mkdir -p lib; cp -f -r jni/* lib; rm -rf jni; done;\
	$$ANDROID_HOME/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/d8 *.classes.jar *.class --release;\
	cp ../AndroidManifest.xml .;\
	$$ANDROID_HOME/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/aapt2 link -o $${APP_ID}.unsigned.apk -I $$ANDROID_HOME/platforms/android-$${ANDROID_API_VERSION}/android.jar --manifest AndroidManifest.xml;\
	zip -ur $${APP_ID}.unsigned.apk classes.dex;\
	zip -ur $${APP_ID}.unsigned.apk lib;\
	mkdir -p ~/.android-certs;\
	[ ! -f ~/.android-certs/$${APP_ID}.keystore ] && keytool -genkeypair -validity 365 -keystore ~/.android-certs/$${APP_ID}.keystore -keyalg RSA -keysize 2048 -keypass $${CERT_PASSWORD} -storepass $${CERT_PASSWORD} -dname CN=$${APP_ID};\
	$$ANDROID_HOME/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/zipalign -f -p 4 $${APP_ID}.unsigned.apk $${APP_ID}.apk;\
	$$ANDROID_HOME/build-tools/$${ANDROID_BUILD_TOOLS_VERSION}/apksigner sign --ks ~/.android-certs/$${APP_ID}.keystore --ks-pass pass:$${CERT_PASSWORD} --key-pass pass:$${CERT_PASSWORD} $${APP_ID}.apk;\

install: build
	$$ANDROID_HOME/platform-tools/adb install build/$${APP_ID}.apk

clean:
	rm -rf build